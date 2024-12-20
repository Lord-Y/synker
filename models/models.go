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
	ValidatedFiles []ValidatedFiles
	// List of validated schemas
	ValidatedSchemas Schemas
	// Init will only perform prerequisites related to elasticsearch / kafka / cockroachdb
	Init bool
}

// List of validated files with SQL queries
type ValidatedFiles struct {
	File    string
	Queries []string
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
	Topic TopicSchema `json:"topic" yaml:"topic" validate:"required"`
	// SQL requirements when advanced query is needed
	SQL SQL `json:"sql" yaml:"sql" validate:"-"`
	// Elasticsearch configuration
	Elasticsearch ElasticsearchSchema `json:"elasticsearch" yaml:"elasticsearch" validate:"required"`
	// // Requirements to create change feed
	ChangeFeed ChangeFeed `json:"changeFeed" yaml:"changeFeed" validate:"required"`
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

// SQL is the requirement to query the SQL database
type SQL struct {
	// ColumnNames is the list of columns name to use to filter kafka message
	// in order to insert or delete data in elasticsearch
	Columns []string `json:"columnNames" yaml:"columnNames" validate:"required"`
	// Query is the SQL query to execute
	Query string `json:"query" yaml:"query" validate:"required"`
}

// ChangeFeed bind all requrirements to create changefeed
type ChangeFeed struct {
	// FullTableName is cockroach full table name like movr.public.promo_codes
	FullTableName string `json:"fullTableName" yaml:"fullTableName" validate:"required"`
	// Options is all change feed options required to create the change feed
	Options []string `json:"options" yaml:"options" validate:"required"`
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
	Partition int `json:"partition"`
	// Offset returned by redpanda/kafka
	Offset int `json:"offset"`
}

// Notify hold all indexes that need to be updated with the appropriate column list
type Notify struct {
	// Indexes list to update
	Indexes []Index `json:"indexes" yaml:"indexes" validate:"required,dive"`
}

// Index hold index name and columns that are used by Notify struct
type Index struct {
	// Name or index alias to use
	Name string `json:"name" yaml:"name" validate:"required"`
	// Column list present in the index to update
	Columns []string `json:"columns" yaml:"columns" validate:"required"`
}
