// Package processing provide all requirements to process change data capture
package processing

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/Lord-Y/synker/elasticsearch"
	"github.com/Lord-Y/synker/kafka"
	"github.com/Lord-Y/synker/logger"
	"github.com/Lord-Y/synker/models"
	"github.com/Lord-Y/synker/tools"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/olivere/elastic/v7"
	"github.com/rs/zerolog/log"
	kafkago "github.com/segmentio/kafka-go"
	"gopkg.in/yaml.v3"
)

type Validate models.Configuration

var (
	// valideFiles will be used when walking into the directory tree
	valideFiles = []string{
		".json",
		".yaml",
		".yml",
	}
	validate *validator.Validate
)

func init() {
	os.Setenv("SYNKER_BATCH_LOG", "true")
	defer os.Unsetenv("APP_BATCH_LOG")
	logger.SetLoggerLogLevel()
}

// Run will validate the files configurations
func (c *Validate) Run() {
	var files []string
	err := filepath.WalkDir(
		c.ConfigDir,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if tools.InSlice(filepath.Ext(d.Name()), valideFiles) {
				files = append(files, path)
			}
			return nil
		},
	)
	if err != nil {
		log.Fatal().Err(err).Msgf("Fail to walk into directory %s", c.ConfigDir)
		return
	}
	if len(files) == 0 {
		log.Fatal().Msg("No config files found to validate")
		return
	}
	c.Files = files
	z, err := c.parsing()
	if err != nil {
		log.Fatal().Err(err).Msgf("File %s is invalid", z)
	}
	for _, v := range c.ValidatedFiles {
		log.Info().Msgf("Config file %s is valid", v.File)
		if len(v.Queries) > 0 {
			for _, q := range v.Queries {
				log.Info().Msgf("SQL queries that will be performed: %s", q)
			}
		}
	}
}

// loadFiles will load all files requests and return a slice of bytes
func loadFiles(f string) (z []byte, err error) {
	of, err := os.Open(f)
	if err != nil {
		return z, fmt.Errorf("Fail to open file `%s` on your system", f)
	}
	defer of.Close()
	fo, err := io.ReadAll(of)
	if err != nil {
		return z, fmt.Errorf("Fail to read file content `%s` on your system", f)
	}
	return fo, err
}

// parsing permit to validate all provided config
func (c *Validate) parsing() (file string, err error) {
	validate = validator.New()

	for _, file = range c.Files {
		fBytes, err := loadFiles(file)
		if err != nil {
			return file, err
		}
		var (
			z       models.Schemas
			zv      models.ValidatedFiles
			queries []string
		)
		if tools.IsYamlFromBytes(fBytes) {
			err = yaml.Unmarshal(fBytes, &z)
			if err != nil {
				return file, err
			}
			err = validate.Struct(z)
			if err != nil {
				return file, err
			}
		}
		if tools.IsJSONFromBytes(fBytes) {
			err = json.Unmarshal(fBytes, &z)
			if err != nil {
				return file, err
			}
			err = validate.Struct(z)
			if err != nil {
				return file, err
			}
		}
		for _, v := range z.Schemas {
			if reflect.ValueOf(v.SQL.QueryType).IsZero() {
				return file, fmt.Errorf("queryType cannot be empty and must be one of advanced, none, notify")
			}
			if !reflect.ValueOf(v.SQL.QueryType.Advanced).IsZero() {
				q := strings.TrimSpace(v.SQL.QueryType.Advanced.Query)
				if q != "" {
					if !strings.Contains(q, ".") {
						return file, fmt.Errorf("Your SQL query `%s` is malformed. It must be in the format SELECT table_name.column_a,table_name.column_b ...", q)
					}
					if strings.Contains(strings.ToLower(q), " where ") && !strings.HasSuffix(q, ")") {
						return file, fmt.Errorf("Your SQL query `%s` is malformed. It has a WHERE condition AND must be in the format SELECT table_name.column_a,table_name.column_b WHERE (table_name.column_c = 1)", q)
					}
					queries = append(queries, q)
				} else {
					return file, fmt.Errorf("SQL query cannot be empty for advanced queryType")
				}
			}
		}
		zv.File = file
		zv.Queries = queries
		c.ValidatedFiles = append(c.ValidatedFiles, zv)
		if len(c.ValidatedSchemas.Schemas) == 0 {
			c.ValidatedSchemas = z
		} else {
			c.ValidatedSchemas.Schemas = append(c.ValidatedSchemas.Schemas, z.Schemas...)
		}
	}
	return
}

// ManageTopics permit to create or update topics
func (c *Validate) ManageTopics() (err error) {
	ls, err := kafka.Client()
	if err != nil {
		return err
	}
	topics, err := kafka.ListTopics(ls)
	if err != nil {
		return err
	}
	for _, v := range c.ValidatedSchemas.Schemas {
		if !tools.InSlice(v.Topic.Name, topics) {
			ct, err := kafka.Client()
			if err != nil {
				return err
			}
			err = kafka.CreateTopic(
				ct,
				models.CreateTopic{
					Name:              v.Topic.Name,
					NumPartitions:     v.Topic.NumPartitions,
					ReplicationFactor: v.Topic.ReplicationFactor,
				},
			)
			if err != nil {
				return err
			}
		}
	}
	return
}

// ManageElasticsearchIndex permit to check or create elasticsearch index
func (c *Validate) ManageElasticsearchIndex() (err error) {
	client, err := elasticsearch.Client()
	if err != nil {
		return err
	}
	defer client.Stop()
	for _, v := range c.ValidatedSchemas.Schemas {
		if v.Elasticsearch.Index.Create {
			ctx := context.Background()
			alias := strings.TrimSpace(v.Elasticsearch.Index.Alias)
			index := strings.TrimSpace(v.Elasticsearch.Index.Name)

			b, err := client.IndexExists(index).Do(ctx)
			if err != nil {
				return err
			}
			if !b {
				create, err := client.CreateIndex(index).
					BodyJson(v.Elasticsearch.Mapping).Do(ctx)
				if err != nil {
					return err
				}
				if !create.Acknowledged {
					return fmt.Errorf("Fail to get index creation ack")
				}
				if alias != "" {
					if alias != "" && index != alias {
						alias_create, err := client.Alias().
							Add(index, alias).
							Action(
								elastic.NewAliasAddAction(alias).
									Index(index).
									IsWriteIndex(true),
							).
							Do(context.TODO())
						if err != nil {
							return err
						}
						if !alias_create.Acknowledged {
							return fmt.Errorf("Fail to get alias creation ack")
						}
					} else {
						return fmt.Errorf("Index %s and alias %s cannot have the same name on schema %s", index, alias, v.Name)
					}
				}
			}
			if alias != "" {
				if alias != "" && index != alias {
					list, err := client.Aliases().
						Index(index).
						Pretty(true).
						Do(context.TODO())
					if err != nil {
						return err
					}
					if len(list.Indices) <= 1 {
						alias_create, err := client.Alias().
							Add(index, alias).
							Action(
								elastic.NewAliasAddAction(alias).
									Index(index).
									IsWriteIndex(true),
							).
							Do(context.TODO())
						if err != nil {
							return err
						}
						if !alias_create.Acknowledged {
							return fmt.Errorf("Fail to get alias creation ack")
						}
					}
				} else {
					return fmt.Errorf("Index %s and alias %s cannot have the same name on schema %s", index, alias, v.Name)
				}
			}
		}
	}
	return
}

// ManageChangeFeed permit check and create required changefeed
func (c *Validate) ManageChangeFeed() (err error) {
	for _, v := range c.ValidatedSchemas.Schemas {
		count, err := countChangeFeed(v.ChangeFeed.FullTableName, "running")
		if err != nil {
			log.Fatal().Err(err).Msgf("Fail to check if required changefeed %s on schema %s has status running", v.ChangeFeed.FullTableName, v.Name)
		}
		if count == 0 {
			err = createChangeFeed(v.ChangeFeed)
			if err != nil {
				log.Fatal().Err(err).Msgf("Fail to create changefeed %s on schema %s", v.ChangeFeed.FullTableName, v.Name)
			}
		}
	}
	return
}

// Processing permit to start processing kafka messages and sent it to elasticsearch
func (c *Validate) Processing() {
	wg := sync.WaitGroup{}

	for k, v := range c.ValidatedSchemas.Schemas {
		wg.Add(1)
		topic := v.Topic.Name
		k := k
		log.Debug().Msgf("Start processing on topic %s", topic)
		go func() {
			defer wg.Done()
			c.consume(k, topic)
		}()
	}
	wg.Wait()
}

// consume permit to consume messages in kafka and sent it to elasticsearch
func (c *Validate) consume(index int, topic string) {
	conn, err := kafka.Client()
	if err != nil {
		return
	}
	defer conn.Close()
	brokers := []string{
		conn.Broker().Host,
		strconv.Itoa(conn.Broker().Port),
	}
	r := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  fmt.Sprintf("synker_%s", strings.TrimSpace(c.ValidatedSchemas.Schemas[index].Name)),
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	defer r.Close()
	if !reflect.ValueOf(c.ValidatedSchemas.Schemas[index].SQL.QueryType.None).IsZero() {
		c.ValidatedSchemas.Schemas[index].SQL.Type = "none"
	}
	if !reflect.ValueOf(c.ValidatedSchemas.Schemas[index].SQL.QueryType.Advanced).IsZero() {
		c.ValidatedSchemas.Schemas[index].SQL.Type = "advanced"
	}

	ctx := context.Background()
	for {
		var (
			message models.ConsumeMessage
			value   map[string]interface{}
			mbytes  []byte
			indexed bool
		)
		m, err := r.FetchMessage(ctx)
		if err != nil {
			break
		}
		log.Debug().Msgf("Message at topic/partition/offset %v/%v/%v: %s = %s", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		mjson, err := json.Marshal(m)
		if err != nil {
			log.Error().Err(err).Msgf("Fail to marshall kafka message in topic %s on partition %d and offset %d", m.Topic, m.Partition, m.Offset)
			return
		}
		err = json.Unmarshal(mjson, &message)
		if err != nil {
			log.Error().Err(err).Msgf("Fail to unmarshall kafka message in topic %s on partition %d and offset %d", m.Topic, m.Partition, m.Offset)
			return
		}
		mbytes, err = base64.StdEncoding.DecodeString(message.Value)
		if err != nil {
			log.Error().Err(err).Msgf("Fail to decode base64 value key from message in topic %s on partition %d and offset %d", m.Topic, m.Partition, m.Offset)
			return
		}
		err = json.Unmarshal(mbytes, &value)
		if err != nil {
			log.Error().Err(err).Msgf("Fail to unmarshall value key from kafka message in topic %s on partition %d and offset %d", m.Topic, m.Partition, m.Offset)
			return
		}
		if value["after"] != nil {
			if c.ValidatedSchemas.Schemas[index].SQL.Type == "none" {
				exist, id, err := c.SearchByVersion(index, value)
				if err != nil {
					log.Error().Err(err).Msgf("Document already exist in elasticsearch with kafka message in topic %s on partition %d and offset %d", m.Topic, m.Partition, m.Offset)
					return
				}
				log.Debug().Msgf("Data exist in elasticsearch? %t", exist)
				var w string
				if exist {
					w = fmt.Sprintf("Fail to index kafka message in topic %s on partition %d and offset %d", m.Topic, m.Partition, m.Offset)
				} else {
					w = fmt.Sprintf("Fail to update elasticsearch index document with kafka message in topic %s on partition %d and offset %d", m.Topic, m.Partition, m.Offset)
				}
				err = c.IndexNewContent(index, value, id)
				if err != nil {
					log.Error().Err(err).Msgf(w)
					return
				}
				indexed = true
			}
			if c.ValidatedSchemas.Schemas[index].SQL.Type == "advanced" {
				t := strings.Split(c.ValidatedSchemas.Schemas[index].ChangeFeed.FullTableName, ".")
				var v map[string]interface{}
				err = mapstructure.Decode(value["after"], &v)
				if err != nil {
					return
				}
				result, select_query, err := c.query(c.ValidatedSchemas.Schemas[index].SQL.Advanced.Query, t[len(t)-1], v)
				if err != nil {
					log.Error().Err(err).Msgf("Fail to execute SQL query `%s` with kafka message in topic %s on partition %d and offset %d", select_query, m.Topic, m.Partition, m.Offset)
					return
				}
				log.Debug().Msgf("result %s query %s", result, select_query)
				exist, id, err := c.SearchByVersion(index, value)
				if err != nil {
					log.Error().Err(err).Msgf("Document already exist in elasticsearch with kafka message in topic %s on partition %d and offset %d", m.Topic, m.Partition, m.Offset)
					return
				}
				log.Debug().Msgf("Data exist in elasticsearch? %t", exist)
				var w string
				if exist {
					w = fmt.Sprintf("Fail to index kafka message in topic %s on partition %d and offset %d", m.Topic, m.Partition, m.Offset)
				} else {
					w = fmt.Sprintf("Fail to update elasticsearch index document with kafka message in topic %s on partition %d and offset %d", m.Topic, m.Partition, m.Offset)
				}
				err = c.IndexNewContent(index, result, id)
				if err != nil {
					log.Error().Err(err).Msgf(w)
					return
				}
				indexed = true
			}
		}
		if indexed {
			log.Debug().Msgf("Kafka message has been indexed in topic %s on partition %d and offset %d", m.Topic, m.Partition, m.Offset)
			if err := r.CommitMessages(ctx, m); err != nil {
				log.Error().Err(err).Msgf("Fail to commit kafka message in topic %s on partition %d and offset %d", m.Topic, m.Partition, m.Offset)
				return
			}
			log.Debug().Msgf("Kafka message with offset %d has been commited in topic %s on partition %d", m.Offset, m.Topic, m.Partition)
		}
	}
}

func (c *Validate) SearchByVersion(index int, value map[string]interface{}) (b bool, id string, err error) {
	var es_target_index string
	client, err := elasticsearch.Client()
	if err != nil {
		return
	}
	defer client.Stop()
	ctx := context.Background()
	es_alias := strings.TrimSpace(c.ValidatedSchemas.Schemas[index].Elasticsearch.Index.Alias)
	es_index := strings.TrimSpace(c.ValidatedSchemas.Schemas[index].Elasticsearch.Index.Name)

	if es_alias != "" {
		es_target_index = es_alias
	} else {
		es_target_index = es_index
	}

	_, err = client.Refresh().Index(es_target_index).Do(ctx)
	if err != nil {
		return
	}

	esQuery := elastic.NewBoolQuery()
	if !reflect.ValueOf(c.ValidatedSchemas.Schemas[index].SQL.None).IsZero() {
		for _, v := range c.ValidatedSchemas.Schemas[index].SQL.None.ImmutableColumns {
			if value["after"] != nil {
				var after map[string]interface{}
				err := mapstructure.Decode(value["after"], &after)
				if err != nil {
					return false, "", err
				}
				if z, ok := after[v.Name]; ok {
					esQuery.Must(elastic.NewMatchQuery(v.Name, z))
				}
			}
		}
	}
	if !reflect.ValueOf(c.ValidatedSchemas.Schemas[index].SQL.Advanced).IsZero() {
		for _, v := range c.ValidatedSchemas.Schemas[index].SQL.Advanced.ImmutableColumns {
			if value["after"] != nil {
				var after map[string]interface{}
				err := mapstructure.Decode(value["after"], &after)
				if err != nil {
					return false, "", err
				}
				if z, ok := after[v.Name]; ok {
					esQuery.Must(elastic.NewMatchQuery(v.Name, z))
				}
			}
		}
	}
	src, err := esQuery.Source()
	if err != nil {
		return
	}
	data, err := json.Marshal(src)
	if err != nil {
		return
	}

	result, err := client.Search().
		Index(es_target_index).
		Query(esQuery).
		FetchSourceContext(elastic.NewFetchSourceContext(true)).
		Do(ctx)
	if err != nil {
		return
	}
	if result.TotalHits() > 0 {
		if len(result.Hits.Hits) > 1 {
			log.Info().Msgf("Elasticsearch search query on index %s %s", es_target_index, string(data))
			return false, "", fmt.Errorf("Multiple document found with same data")
		}
		for _, hit := range result.Hits.Hits {
			id = hit.Id
		}
		return true, id, nil
	}
	return
}

func (c *Validate) IndexNewContent(index int, value map[string]interface{}, id string) (err error) {
	var es_target_index string
	client, err := elasticsearch.Client()
	defer client.Stop()

	es_alias := strings.TrimSpace(c.ValidatedSchemas.Schemas[index].Elasticsearch.Index.Alias)
	es_index := strings.TrimSpace(c.ValidatedSchemas.Schemas[index].Elasticsearch.Index.Name)

	var content map[string]interface{}
	switch c.ValidatedSchemas.Schemas[index].SQL.Type {
	case "none":
		err = mapstructure.Decode(value["after"], &content)
		if err != nil {
			return
		}
	case "advanced":
		content = value
	}

	if es_alias != "" {
		es_target_index = es_alias
	} else {
		es_target_index = es_index
	}
	ctx := context.Background()
	if id == "" {
		uuidgen := uuid.New()
		_, err = client.
			Index().
			Index(es_target_index).
			Id(uuidgen.String()).
			BodyJson(content).
			Do(ctx)
	} else {
		_, err = client.Update().
			Index(es_target_index).
			Id(id).
			Doc(content).
			Do(ctx)
	}
	if err != nil {
		return
	}
	return
}
