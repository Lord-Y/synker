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

	"github.com/Lord-Y/synker/tools"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/nqd/flat"
	"github.com/olivere/elastic/v7"
	kafkago "github.com/segmentio/kafka-go"
	"gopkg.in/yaml.v3"
)

type Validate Configuration

var (
	// valideFiles will be used when walking into the directory tree
	valideFiles = []string{
		".json",
		".yaml",
		".yml",
	}
	validate *validator.Validate
)

// ParseAndValidateConfig will validate the files configurations
func (c *Validate) ParseAndValidateConfig() {
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
		c.Logger.Fatal().Err(err).Msgf("Fail to walk into directory %s", c.ConfigDir)
		return
	}
	if len(files) == 0 {
		c.Logger.Fatal().Msg("No config files found to validate")
		return
	}
	c.files = files
	z, err := c.parsing()
	if err != nil {
		c.Logger.Fatal().Err(err).Msgf("File %s is invalid", z)
	}
	for _, v := range c.validatedFiles {
		c.Logger.Info().Msgf("Config file %s is valid", v.File)
		if len(v.Queries) > 0 {
			for _, query := range v.Queries {
				c.Logger.Info().Msgf("SQL query that will be executed: %s", query)
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

	for _, file = range c.files {
		fBytes, err := loadFiles(file)
		if err != nil {
			return file, err
		}
		var (
			z       schemas
			zv      validatedFiles
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
			if !reflect.ValueOf(v.SQL).IsZero() {
				query := strings.TrimSpace(v.SQL.Query)
				if query != "" {
					if !strings.Contains(query, ".") {
						return file, fmt.Errorf("Your SQL query `%s` is malformed. It must be in the format SELECT table_name.column_a,table_name.column_b ...", query)
					}
					if strings.Contains(strings.ToLower(query), " where ") && !strings.HasSuffix(query, ")") {
						return file, fmt.Errorf("Your SQL query `%s` is malformed. It has a WHERE condition AND must be in the format SELECT table_name.column_a,table_name.column_b WHERE (table_name.column_c = 1)", query)
					}
					queries = append(queries, query)
				} else {
					return file, fmt.Errorf("SQL query cannot be empty for advanced queryType")
				}
			}
		}
		zv.File = file
		zv.Queries = queries
		c.validatedFiles = append(c.validatedFiles, zv)
		if len(c.validatedSchemas.Schemas) == 0 {
			c.validatedSchemas = z
		} else {
			c.validatedSchemas.Schemas = append(c.validatedSchemas.Schemas, z.Schemas...)
		}
	}
	return
}

// manageTopics permit to create or update topics
func (c *Validate) manageTopics() (err error) {
	ls, err := c.kClient()
	if err != nil {
		return err
	}

	topics, err := c.listTopics(ls)
	if err != nil {
		return err
	}

	for _, v := range c.validatedSchemas.Schemas {
		if !tools.InSlice(v.Topic.Name, topics) {
			client, err := c.kClient()
			if err != nil {
				return err
			}
			err = c.createTopic(
				client,
				createTopicModel{
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

// manageElasticsearchIndex permit to check or create elasticsearch index
func (c *Validate) manageElasticsearchIndex() (err error) {
	client, err := c.eClient()
	if err != nil {
		return err
	}
	defer client.Stop()
	for _, v := range c.validatedSchemas.Schemas {
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
					return fmt.Errorf("Fail to get index creation acknowledgement")
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
							return fmt.Errorf("Fail to get alias creation acknowledgement")
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
							return fmt.Errorf("Fail to get alias creation acknowledgement")
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

// manageChangeFeed permit check and create required changefeed
func (c *Validate) manageChangeFeed() (err error) {
	for _, v := range c.validatedSchemas.Schemas {
		count, err := countChangeFeed(v.ChangeFeed.FullTableName, "running")
		if err != nil {
			c.Logger.Fatal().Err(err).Msgf("Fail to check if required changefeed %s on schema %s has status running", v.ChangeFeed.FullTableName, v.Name)
		}
		if count == 0 {
			err = createChangeFeed(v.ChangeFeed)
			if err != nil {
				c.Logger.Fatal().Err(err).Msgf("Fail to create changefeed %s on schema %s", v.ChangeFeed.FullTableName, v.Name)
			}
		}
	}
	return
}

// processing permit to start processing kafka messages and sent it to elasticsearch
func (c *Validate) processing() {
	wg := sync.WaitGroup{}

	for k, v := range c.validatedSchemas.Schemas {
		wg.Add(1)
		topic := v.Topic.Name
		k := k
		c.Logger.Debug().Msgf("Start processing on topic %s", topic)
		go func() {
			defer wg.Done()
			c.consume(k, topic)
		}()
	}
	wg.Wait()
}

// consume permit to consume messages in kafka and sent it to elasticsearch
func (c *Validate) consume(index int, topic string) {
	conn, err := c.kClient()
	if err != nil {
		return
	}
	defer conn.Close()

	connLeader, err := c.connectToController(conn)
	if err != nil {
		c.Logger.Fatal().Err(err).Msg("Fail to connect to the controller")
	}
	defer connLeader.Close()

	brokerList, err := connLeader.Brokers()
	if err != nil {
		c.Logger.Fatal().Err(err).Msg("Brokers list cannot be empty")
	}

	var brokers []string
	for _, v := range brokerList {
		brokers = append(brokers, v.Host+":"+strconv.Itoa(conn.Broker().Port))
	}

	r := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  fmt.Sprintf("synker_%s", strings.TrimSpace(c.validatedSchemas.Schemas[index].Name)),
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	defer r.Close()

	ctx := context.Background()
	for {
		var (
			message                          consumeMessage
			value                            map[string]interface{}
			key                              []string
			mkBytes, mvBytes                 []byte
			documentIndexed, documentDeleted bool
			esIndex                          string
		)

		m, err := r.FetchMessage(ctx)
		if err != nil {
			c.Logger.Fatal().Err(err).Msgf("Fail to fetch kafka message from topic %s on partition %d and offset %d", topic, m.Partition, m.Offset)
		}

		c.Logger.Debug().Msgf("Message at topic/partition/offset %v/%v/%v: %s = %s", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		mjson, err := json.Marshal(m)
		if err != nil {
			c.increaseMetrics("kafka", m.Topic, "marshalling")
			c.Logger.Error().Err(err).Msgf("Fail to marshal kafka message from topic %s on partition %d and offset %d", m.Topic, m.Partition, m.Offset)
			return
		}

		err = json.Unmarshal(mjson, &message)
		if err != nil {
			c.increaseMetrics("kafka", m.Topic, "marshalling")
			c.Logger.Error().Err(err).Msgf("Fail to unmarshal kafka message from topic %s on partition %d and offset %d", m.Topic, m.Partition, m.Offset)
			return
		}

		mkBytes, err = base64.StdEncoding.DecodeString(message.Key)
		if err != nil {
			c.increaseMetrics("kafka", m.Topic, "encoding")
			c.Logger.Error().Err(err).Msgf("Fail to decode base64 field value from message from topic %s on partition %d and offset %d", m.Topic, m.Partition, m.Offset)
			return
		}

		err = json.Unmarshal(mkBytes, &key)
		if err != nil {
			c.increaseMetrics("kafka", m.Topic, "marshalling")
			c.Logger.Error().Err(err).Msgf("Fail to unmarshal field key from kafka message from topic %s on partition %d and offset %d", m.Topic, m.Partition, m.Offset)
			return
		}

		mvBytes, err = base64.StdEncoding.DecodeString(message.Value)
		if err != nil {
			c.increaseMetrics("kafka", m.Topic, "encoding")
			c.Logger.Error().Err(err).Msgf("Fail to decode base64 field value from message from topic %s on partition %d and offset %d", m.Topic, m.Partition, m.Offset)
			return
		}

		err = json.Unmarshal(mvBytes, &value)
		if err != nil {
			c.increaseMetrics("kafka", m.Topic, "marshalling")
			c.Logger.Error().Err(err).Msgf("Fail to unmarshal field value from kafka message from topic %s on partition %d and offset %d", m.Topic, m.Partition, m.Offset)
			return
		}

		if value["after"] != nil {
			if reflect.ValueOf(c.validatedSchemas.Schemas[index].SQL).IsZero() {
				exist, id, esTargetIndex, err := c.searchByVersion(index, value)
				if err != nil {
					c.increaseMetrics("elasticsearch", m.Topic, "indexing")
					c.Logger.Error().Err(err).Msgf("Document already exist in elasticsearch index %s with kafka message from topic %s on partition %d and offset %d", esTargetIndex, m.Topic, m.Partition, m.Offset)
					return
				}

				c.Logger.Debug().Msgf("Data exist in elasticsearch index %s? %t", esTargetIndex, exist)
				err = c.indexNewContent(index, value, id)
				if err != nil {
					c.increaseMetrics("elasticsearch", m.Topic, "indexing")
					c.Logger.Error().Err(err).Msgf("Fail to update elasticsearch document in index %s with kafka message from topic %s on partition %d and offset %d", esTargetIndex, m.Topic, m.Partition, m.Offset)
					return
				}
				documentIndexed = true
				esIndex = esTargetIndex
			} else {
				t := strings.Split(c.validatedSchemas.Schemas[index].ChangeFeed.FullTableName, ".")
				var v map[string]interface{}
				err = mapstructure.Decode(value["after"], &v)
				if err != nil {
					return
				}

				result, select_query, err := c.query(c.validatedSchemas.Schemas[index].SQL.Query, t[len(t)-1], v)
				if err != nil {
					c.increaseMetrics("elasticsearch", m.Topic, "sql")
					c.Logger.Error().Err(err).Msgf("Fail to execute SQL query `%s` with kafka message from topic %s on partition %d and offset %d", select_query, m.Topic, m.Partition, m.Offset)
					return
				}

				c.Logger.Debug().Msgf("result %+v query %+v", result, select_query)
				exist, id, esTargetIndex, err := c.searchByVersion(index, value)
				if err != nil {
					c.increaseMetrics("elasticsearch", m.Topic, "indexing")
					c.Logger.Error().Err(err).Msgf("Document already exist in elasticsearch index %s with kafka message from topic %s on partition %d and offset %d", esTargetIndex, m.Topic, m.Partition, m.Offset)
					return
				}

				c.Logger.Debug().Msgf("Data exist in elasticsearch? %t", exist)
				var errorMessage string
				if exist {
					errorMessage = fmt.Sprintf("Fail to index kafka message from topic %s on partition %d and offset %d", m.Topic, m.Partition, m.Offset)
				} else {
					errorMessage = fmt.Sprintf("Fail to update elasticsearch document in index %s with kafka message from topic %s on partition %d and offset %d", esTargetIndex, m.Topic, m.Partition, m.Offset)
				}

				err = c.indexNewContent(index, result, id)
				if err != nil {
					c.increaseMetrics("elasticsearch", m.Topic, "indexing")
					c.Logger.Error().Err(err).Msgf("%s", errorMessage)
					return
				}
				documentIndexed = true
				esIndex = esTargetIndex
			}

			if documentIndexed {
				c.Logger.Debug().Msgf("Kafka message has been indexed into elasticsearch index `%s` from topic %s on partition %d and offset %d", esIndex, m.Topic, m.Partition, m.Offset)
				if err := r.CommitMessages(ctx, m); err != nil {
					c.increaseMetrics("kafka", m.Topic, "commit")
					c.Logger.Error().Err(err).Msgf("Fail to commit kafka message from topic %s on partition %d and offset %d", m.Topic, m.Partition, m.Offset)
					return
				}
				c.Logger.Debug().Msgf("Kafka message has been commited in topic %s on partition %d and offset %d", m.Topic, m.Partition, m.Offset)
			}
		} else {
			exist, id, esTargetIndex, err := c.searchByVersion(index, value)
			if err != nil {
				c.increaseMetrics("elasticsearch", m.Topic, "indexing")
				c.Logger.Error().Err(err).Msgf("Document with key(s) %s in elasticsearch index %s with kafka message from topic %s on partition %d and offset %d", key, esTargetIndex, m.Topic, m.Partition, m.Offset)
				return
			}

			c.Logger.Debug().Msgf("Data exist in elasticsearch index %s? %t", esTargetIndex, exist)

			if exist {
				err = c.deleteContent(index, id)
				if err != nil {
					c.increaseMetrics("elasticsearch", m.Topic, "delete")
					c.Logger.Error().Err(err).Msgf("Fail to delete elasticsearch document with id %s in index %s with kafka message from topic %s on partition %d and offset %d", id, esTargetIndex, m.Topic, m.Partition, m.Offset)
					return
				}
				documentDeleted = true
				esIndex = esTargetIndex
			}

			if documentDeleted {
				c.Logger.Debug().Msgf("Kafka message has been deleted from elasticsearch index `%s` from topic %s on partition %d and offset %d", esIndex, m.Topic, m.Partition, m.Offset)
				if err := r.CommitMessages(ctx, m); err != nil {
					c.increaseMetrics("kafka", m.Topic, "commit")
					c.Logger.Error().Err(err).Msgf("Fail to commit kafka message from topic %s on partition %d and offset %d", m.Topic, m.Partition, m.Offset)
					return
				}
				c.Logger.Debug().Msgf("Kafka message has been commited in topic %s on partition %d and offset %d", m.Topic, m.Partition, m.Offset)
			}
		}
	}
}

// searchByVersion permit to search into elasticsearch
//
// if the new data provided already exist
//
// if it exist, it will return the elasticsearch id
//
// if it exist multiple times, an error will be returned
func (c *Validate) searchByVersion(index int, value map[string]interface{}) (found bool, id, esTargetIndex string, err error) {
	client, err := c.eClient()
	if err != nil {
		return
	}
	defer client.Stop()

	ctx := context.Background()
	esAlias := strings.TrimSpace(c.validatedSchemas.Schemas[index].Elasticsearch.Index.Alias)
	esIndex := strings.TrimSpace(c.validatedSchemas.Schemas[index].Elasticsearch.Index.Name)

	if esAlias != "" {
		esTargetIndex = esAlias
	} else {
		esTargetIndex = esIndex
	}

	_, err = client.Refresh().Index(esTargetIndex).Do(ctx)
	if err != nil {
		return
	}

	var documentToDelete bool
	esQuery := elastic.NewBoolQuery()
	if reflect.ValueOf(c.validatedSchemas.Schemas[index].SQL).IsZero() {
		var queries []elastic.Query
		if value["after"] != nil {
			var before map[string]interface{}
			if value["before"] != nil {
				err := mapstructure.Decode(value["before"], &before)
				if err != nil {
					return false, "", esTargetIndex, err
				}
				for k, v := range before {
					isMap, multiKeys := c.isMap(v)
					switch {
					case isMap && multiKeys:
						continue
					case isMap && !multiKeys:
						x := make(map[string]interface{})
						x[k] = v
						flatten, err := flat.Flatten(x, nil)
						if err != nil {
							c.Logger.Error().Err(err).Msgf("Fail to flatten map %+v", v)
						} else {
							c.Logger.Debug().Msgf("flatten map %+v", flatten)
							for fk, fv := range flatten {
								c.Logger.Debug().Msgf("K/V to search for %v == %v", fk, fv)
								queries = append(queries, elastic.NewMatchQuery(fk, fv))
							}
						}
					default:
						c.Logger.Debug().Msgf("K/V to search for %v == %v", k, v)
						queries = append(queries, elastic.NewMatchQuery(k, v))
					}
				}
			} else {
				return
			}
		} else {
			documentToDelete = true

			var before map[string]interface{}
			err := mapstructure.Decode(value["before"], &before)
			if err != nil {
				return false, "", esTargetIndex, err
			}
			for k, v := range before {
				isMap, multiKeys := c.isMap(v)
				switch {
				case isMap && multiKeys:
					continue
				case isMap && !multiKeys:
					x := make(map[string]interface{})
					x[k] = v
					flatten, err := flat.Flatten(x, nil)
					if err != nil {
						c.Logger.Error().Err(err).Msgf("Fail to flatten map %+v", v)
					} else {
						c.Logger.Debug().Msgf("flatten map %+v", flatten)
						for fk, fv := range flatten {
							c.Logger.Debug().Msgf("K/V to delete for %v == %v", fk, fv)
							queries = append(queries, elastic.NewMatchQuery(fk, fv))
						}
					}
				default:
					c.Logger.Debug().Msgf("K/V to delete for %v == %v", k, v)
					queries = append(queries, elastic.NewMatchQuery(k, v))
				}
			}
		}
		esQuery.Must(queries...)
	} else {
		var queries []elastic.Query
		for _, column := range c.validatedSchemas.Schemas[index].SQL.Columns {
			if value["after"] != nil {
				if value["before"] != nil {
					var before map[string]interface{}
					err := mapstructure.Decode(value["before"], &before)
					if err != nil {
						return false, "", esTargetIndex, err
					}
					if z, ok := before[column]; ok {
						c.Logger.Debug().Msgf("K/V to search for %v == %v", column, z)
						queries = append(queries, elastic.NewMatchQuery(column, z))
					}
				} else {
					return
				}
			} else {
				documentToDelete = true

				var before map[string]interface{}
				err := mapstructure.Decode(value["before"], &before)
				if err != nil {
					return false, "", esTargetIndex, err
				}
				if z, ok := before[column]; ok {
					c.Logger.Debug().Msgf("K/V to delete %v == %v", column, z)
					queries = append(queries, elastic.NewMatchQuery(column, z))
				}
			}
		}
		esQuery.Must(queries...)
	}

	src, err := esQuery.Source()
	if err != nil {
		return
	}

	data, err := json.Marshal(src)
	if err != nil {
		return
	}

	c.Logger.Debug().Msgf("Elasticsearch query that will be executed %s", string(data))
	result, err := client.Search().
		Index(esTargetIndex).
		Query(esQuery).
		FetchSourceContext(elastic.NewFetchSourceContext(true)).
		Do(ctx)
	if err != nil {
		return
	}

	if result.TotalHits() > 0 {
		if documentToDelete {
			for _, hit := range result.Hits.Hits {
				id = hit.Id
			}
			return true, id, esTargetIndex, nil
		}

		if len(result.Hits.Hits) > 1 {
			c.Logger.Info().Msgf("Elasticsearch search query on index %s %s", esTargetIndex, string(data))
			return false, "", esTargetIndex, fmt.Errorf("Multiple document found with same data")
		}
		for _, hit := range result.Hits.Hits {
			id = hit.Id
		}
		return true, id, esTargetIndex, nil
	}
	return
}

// indexNewContent permit to add or update provided data
// into elasticsearch index
func (c *Validate) indexNewContent(index int, value map[string]interface{}, uniqId string) (err error) {
	var esTargetIndex string
	client, err := c.eClient()
	defer client.Stop()

	esAlias := strings.TrimSpace(c.validatedSchemas.Schemas[index].Elasticsearch.Index.Alias)
	esIndex := strings.TrimSpace(c.validatedSchemas.Schemas[index].Elasticsearch.Index.Name)

	var content map[string]interface{}
	if !reflect.ValueOf(c.validatedSchemas.Schemas[index].SQL).IsZero() {
		content = value
	} else {
		err = mapstructure.Decode(value["after"], &content)
		if err != nil {
			return
		}
	}

	if esAlias != "" {
		esTargetIndex = esAlias
	} else {
		esTargetIndex = esIndex
	}

	ctx := context.Background()
	if uniqId == "" {
		uuidgen := uuid.New()
		_, err = client.
			Index().
			Index(esTargetIndex).
			Id(uuidgen.String()).
			BodyJson(content).
			Do(ctx)
	} else {
		_, err = client.Update().
			Index(esTargetIndex).
			Id(uniqId).
			Doc(content).
			Do(ctx)
	}

	if err != nil {
		return
	}
	return
}

// deleteContent permit to delete data with the provided id
// from elasticsearch index
func (c *Validate) deleteContent(index int, id string) (err error) {
	var esTargetIndex string
	client, err := c.eClient()
	defer client.Stop()

	esAlias := strings.TrimSpace(c.validatedSchemas.Schemas[index].Elasticsearch.Index.Alias)
	esIndex := strings.TrimSpace(c.validatedSchemas.Schemas[index].Elasticsearch.Index.Name)

	if esAlias != "" {
		esTargetIndex = esAlias
	} else {
		esTargetIndex = esIndex
	}

	ctx := context.Background()

	_, err = client.
		Delete().
		Index(esTargetIndex).
		Id(id).
		Do(ctx)

	if err != nil {
		return
	}
	return
}

func (c *Validate) isMap(x interface{}) (isMap bool, multiKeys bool) {
	t := fmt.Sprintf("%T", x)
	var m map[string]interface{}
	isMap = strings.HasPrefix(t, "map[")
	if isMap {
		err := mapstructure.Decode(x, &m)
		if err != nil {
			c.Logger.Error().Err(err).Msgf("Fail to decode map structure")
			return isMap, false
		}
		if len(m) > 1 {
			return true, true
		}
	}
	return
}
