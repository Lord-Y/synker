// Package elasticsearch permit to push or delete all datas from db to elasticsearch
package elasticsearch

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	assert := assert.New(t)
	b := Ping()
	assert.Equal(true, b)
}

func TestPing_fail_auth(t *testing.T) {
	assert := assert.New(t)
	b := Ping()
	assert.Equal(true, b)
}

func TestPing_fail(t *testing.T) {
	assert := assert.New(t)

	os.Setenv("SYNKER_ELASTICSEARCH_USER", "youhou")
	os.Setenv("SYNKER_ELASTICSEARCH_PASSWORD", "youhou")
	defer os.Unsetenv("SYNKER_ELASTICSEARCH_USER")
	defer os.Unsetenv("SYNKER_ELASTICSEARCH_PASSWORD")
	b := Ping()
	assert.Equal(true, b)
}

func TestClient(t *testing.T) {
	assert := assert.New(t)

	client, err := Client()
	defer client.Stop()
	assert.Nil(err)
}

func TestCreateIndex(t *testing.T) {
	assert := assert.New(t)

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

	client, err := Client()
	assert.Nil(err)

	b, err := createIndex(client, "twitter", mapping)
	assert.Nil(err)
	assert.Equal(true, b)
}

func TestCreateIndex_fail(t *testing.T) {
	assert := assert.New(t)

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

	os.Setenv("SYNKER_ELASTICSEARCH_USER", "youhou")
	os.Setenv("SYNKER_ELASTICSEARCH_PASSWORD", "youhou")
	defer os.Unsetenv("SYNKER_ELASTICSEARCH_USER")
	defer os.Unsetenv("SYNKER_ELASTICSEARCH_PASSWORD")

	client, err := Client()
	assert.Nil(err)

	b, err := createIndex(client, "twitter", mapping)
	assert.Nil(err)
	assert.Equal(true, b)
}

func TestIndexAlreadyExist(t *testing.T) {
	assert := assert.New(t)

	client, err := Client()
	assert.Nil(err)

	b, err := IndexAlreadyExist(client, "twitter")
	assert.Nil(err)
	assert.Equal(true, b)
}

func TestDeleteIndex(t *testing.T) {
	assert := assert.New(t)

	client, err := Client()
	assert.Nil(err)

	b, err := DeleteIndex(client, "twitter")
	assert.Nil(err)
	assert.Equal(true, b)
}
