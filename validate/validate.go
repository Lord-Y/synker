// Package validate permit to validate files provided
package validate

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Lord-Y/synker/elasticsearch"
	"github.com/Lord-Y/synker/kafka"
	"github.com/Lord-Y/synker/logger"
	"github.com/Lord-Y/synker/models"
	"github.com/Lord-Y/synker/tools"
	"github.com/olivere/elastic/v7"
	"github.com/rs/zerolog/log"
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
)

func init() {
	logger.SetCLILoggerLogLevel()
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
		log.Error().Err(err).Msgf("Fail to walk into directory %s", c.ConfigDir)
		return
	}
	if len(files) == 0 {
		log.Fatal().Msg("No config files found to validate")
		return
	}
	c.Files = files
	z, err := c.parsing()
	if err != nil {
		log.Error().Err(err).Msgf("File %s is invalid", z)
		return
	}
	for _, v := range c.ValidatedFiles {
		log.Info().Msgf("Config file %s is valid", v)
	}
}

// loadFiles will load all files requests and return a slice of bytes
func loadFiles(f string) (z []byte, err error) {
	of, err := os.Open(f)
	if err != nil {
		return z, fmt.Errorf("Fail to open file `%s` on your system", f)
	}
	defer of.Close()
	fo, err := ioutil.ReadAll(of)
	if err != nil {
		return z, fmt.Errorf("Fail to read file content `%s` on your system", f)
	}
	return fo, err
}

// parsing permit to validate all provided config
func (c *Validate) parsing() (file string, err error) {
	for _, file = range c.Files {
		fBytes, err := loadFiles(file)
		if err != nil {
			return file, err
		}
		var z models.Schemas
		if tools.IsYamlFromBytes(fBytes) {
			err = yaml.Unmarshal(fBytes, &z)
			if err != nil {
				return file, err
			}
		}
		if tools.IsJSONFromBytes(fBytes) {
			err = json.Unmarshal(fBytes, &z)
			if err != nil {
				return file, err
			}
		}
		c.ValidatedFiles = append(c.ValidatedFiles, file)
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
					NumPartitions:     1,
					ReplicationFactor: 3,
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
