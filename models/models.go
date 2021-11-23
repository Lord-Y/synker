// Package models assemble all structs, interface e.g ...
package models

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
	Name string `json:"name" yaml:"name"`
	// Number of partitions
	NumPartitions int `json:"numPartitions" yaml:"numPartitions"`
	// Cluster replication factor
	ReplicationFactor int `json:"replicationFactor" yaml:"replicationFactor"`
	// Topic config
	TopicConfig []TopicConfig
}

// TopicConfig store key/value pair use to configure the topic
type TopicConfig struct {
	Key   string `json:"key" yaml:"key"`     // Config key
	Value string `json:"value" yaml:"value"` // Config value
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
	Schemas []ConfigSchema `json:"schemas" yaml:"schemas" binding:"required"`
}

// ConfigSchema is the validator
type ConfigSchema struct {
	// Schema name
	Name string `json:"name" yaml:"name" binding:"required"`
	// Topic schema
	Topic TopicSchema `json:"topic" yaml:"topic"`
	// SQL schema
	SQL SQLSchema `json:"sql" yaml:"sql"`
	// Elasticsearch configuration
	Elasticsearch ElasticsearchSchema `json:"elasticsearch" yaml:"elasticsearch"`
}

// TopicSchema is the requirement to create the topic
type TopicSchema struct {
	// Topic name
	Name string `json:"name" yaml:"name" binding:"required"`
	// Number of partitions
	NumPartitions int `json:"numPartitions" yaml:"numPartitions"`
	// Cluster replication factor
	ReplicationFactor int `json:"replicationFactor" yaml:"replicationFactor"`
	// Topic config
	TopicConfig []TopicConfig `json:"config" yaml:"config"`
}

// SQLSchema is the requirement to query the SQL database
type SQLSchema struct {
	// Type defined if the sql query is plain or not
	Type string `json:"type" yaml:"type" binding:"required,oneof=plain"`
	// Plain query
	Query string `json:"query" yaml:"query"`
}

// ElasticsearchSchema is the requirement related to elasticsearch
type ElasticsearchSchema struct {
	// Elasticsearch index
	Index ElasticsearchIndex `json:"index" yaml:"index"`
	// Type defined if the sql query is plain or not
	Mapping map[string]interface{} `json:"mapping" yaml:"mapping"`
}

// ElasticsearchIndex is the requirement to the elasticsearch index
type ElasticsearchIndex struct {
	// Create index
	Create bool `json:"create" yaml:"create"`
	// Index name
	Name string `json:"name" yaml:"name"`
	// Alias name
	Alias string `json:"alias" yaml:"alias"`
}
