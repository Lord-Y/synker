// Package processing provide all requirements to process change data capture
package processing

import (
	"os"
	"testing"

	"github.com/Lord-Y/synker/logger"
	"github.com/stretchr/testify/assert"
)

func TestElasticsearchPing(t *testing.T) {
	assert := assert.New(t)
	var c Validate
	c.Logger = logger.NewLogger()

	b := c.ePing()
	assert.Equal(true, b)
}

func TestElasticsearchPing_fail(t *testing.T) {
	assert := assert.New(t)
	var c Validate
	c.Logger = logger.NewLogger()

	os.Setenv("SYNKER_ELASTICSEARCH_URI", "http://127.0.0.1:19200")
	b := c.ePing()
	assert.Equal(false, b)
	os.Setenv("SYNKER_ELASTICSEARCH_URI", "http://127.0.0.1:9200")
}

func TestElasticsearchClient(t *testing.T) {
	assert := assert.New(t)
	var c Validate
	c.Logger = logger.NewLogger()

	client, err := c.eClient()
	defer client.Stop()
	assert.Nil(err)
}

func TestElasticsearchCreateIndex(t *testing.T) {
	assert := assert.New(t)
	var c Validate
	c.Logger = logger.NewLogger()

	mapping := `
	{
		"settings":{
			"number_of_shards":1,
			"number_of_replicas":0
		},
		"mappings":{
			"properties":{
				"user":{
					"type":"keyword"
				},
				"message":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
				"retweets":{
					"type":"long"
				},
				"tags":{
					"type":"keyword"
				},
				"location":{
					"type":"geo_point"
				},
				"suggest_field":{
					"type":"completion"
				}
			}
		}
	}
	`

	client, err := c.eClient()
	assert.Nil(err)

	b, err := c.createIndex(client, "twitter", mapping)
	assert.Nil(err)
	assert.Equal(true, b)
}

func TestElasticsearchIndexAlreadyExist(t *testing.T) {
	assert := assert.New(t)
	var c Validate
	c.Logger = logger.NewLogger()

	client, err := c.eClient()
	assert.Nil(err)

	b, err := c.indexAlreadyExist(client, "twitter")
	assert.Nil(err)
	assert.Equal(true, b)
}

func TestElasticsearchDeleteIndex(t *testing.T) {
	assert := assert.New(t)
	var c Validate
	c.Logger = logger.NewLogger()

	client, err := c.eClient()
	assert.Nil(err)

	b, err := c.deleteIndex(client, "twitter")
	assert.Nil(err)
	assert.Equal(true, b)
}
