// Package models assemble all structs, interface e.g ...
package models

import "time"

// Configuration reference all requirements relation to the file config
type Configuration struct {
	// Config dir containing files to be parsed
	ConfigDir string `json:"configDir" yaml:"configDir"`
	// List of files returned after walking into specified directory
	Files []string
	// List of validated files
	ValidatedFiles []string
	// List of validated schemas
	ValidatedSchemas Schemas
}

// CreateTopic reference all the possible requirements to create a topic
type CreateTopic struct {
	// Topic name
	Name string `json:"name" yaml:"name" validate:"required"`
	// Number of partitions
	NumPartitions int `json:"numPartitions" yaml:"numPartitions" validate:"required"`
	// Cluster replication factor
	ReplicationFactor int `json:"replicationFactor" yaml:"replicationFactor" validate:"required"`
	// Topic config
	TopicConfig []TopicConfig
}

// TopicConfig store key/value pair use to configure the topic
type TopicConfig struct {
	Key   string `json:"key" yaml:"key" validate:"required"`     // Config key
	Value string `json:"value" yaml:"value" validate:"required"` // Config value
}

// KafkaWriteMessage define the requirements to write messages into kafka
type KafkaWriteMessage struct {
	TopicName string `json:"name" yaml:"name"`   // Topic name
	Key       string `json:"key" yaml:"key"`     // Message key
	Value     string `json:"value" yaml:"value"` // Message value
}

// Schemas represent the global config to push data to elasticsearch
type Schemas struct {
	// Config schema
	Schemas []ConfigSchema `json:"schemas" yaml:"schemas" validate:"required,dive"`
}

// ConfigSchema is the validator
type ConfigSchema struct {
	// Schema name
	Name string `json:"name" yaml:"name" validate:"required"`
	// Topic schema
	Topic TopicSchema `json:"topic" yaml:"topic" validate:"required,dive"`
	// SQL schema
	SQL SQLSchema `json:"sql" yaml:"sql" validate:"required,dive"`
	// Elasticsearch configuration
	Elasticsearch ElasticsearchSchema `json:"elasticsearch" yaml:"elasticsearch" validate:"required,dive"`
}

// TopicSchema is the requirement to create the topic
type TopicSchema struct {
	// Topic name
	Name string `json:"name" yaml:"name" validate:"required"`
	// Number of partitions
	NumPartitions int `json:"numPartitions" yaml:"numPartitions" validate:"required"`
	// Cluster replication factor
	ReplicationFactor int `json:"replicationFactor" yaml:"replicationFactor" validate:"required"`
	// Topic config
	TopicConfig []TopicConfig `json:"config" yaml:"config" validate:"dive"`
}

// SQLSchema is the requirement to query the SQL database
type SQLSchema struct {
	// Type defined if the sql query is simple or not
	Type string `json:"type" yaml:"type" validate:"required,oneof=simple"`
	// Simple query
	Query string `json:"query" yaml:"query" validate:"required"`
	// Immutable columns are used in where clause to query cockroach or elasticsearch, they can be primary keys e.g
	ImmutableColumns []ImmutableColumn `json:"immutableColumns" yaml:"immutableColumns" validate:"required,dive"`
}

// ImmutableColumn bind sql column name and type
type ImmutableColumn struct {
	// Column name
	Name string `json:"name" yaml:"name" validate:"required"`
	// Column type
	Type string `json:"type" yaml:"type" validate:"required,oneof=integer uuid varchar date"`
}

// ElasticsearchSchema is the requirement related to elasticsearch
type ElasticsearchSchema struct {
	// Elasticsearch index
	Index ElasticsearchIndex `json:"index" yaml:"index" validate:"required"`
	// Type defined if the sql query is plain or not
	Mapping map[string]interface{} `json:"mapping" yaml:"mapping" validate:"required"`
}

// ElasticsearchIndex is the requirement to the elasticsearch index
type ElasticsearchIndex struct {
	// Create index
	Create bool `json:"create" yaml:"create"`
	// Index name
	Name string `json:"name" yaml:"name" validate:"required"`
	// Alias name
	Alias string `json:"alias" yaml:"alias"`
}

// ConsumeMessage permit to consume kafka messages
type ConsumeMessage struct {
	// Topic name returned by cockroach
	Topic string `json:"topic"`
	// Key returned by cockroach
	Key string `json:"key"`
	// Value is the mapping return by cockroach
	Value string `json:"value"`
	// Timestamp returned by cockroach
	Timestamp time.Time `json:"timestamp"`
	// Updated timestamp returned by cockroach
	Updated time.Time `json:"updated"`
	// Partition returned by redpanda/kafka
	Partition int `json:"parition"`
	// Offset returned by redpanda/kafka
	Offset int `json:"offset"`
}
