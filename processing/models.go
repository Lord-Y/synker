// Package processing provide all requirements to process change data capture
package processing

import (
	"time"

	"github.com/rs/zerolog"
)

// Configuration reference all requirements relation to the file config
type Configuration struct {
	// Config dir containing files to be parsed
	ConfigDir string `json:"configDir" yaml:"configDir"`
	// List of files returned after walking into specified directory
	files []string
	// List of validated files
	validatedFiles []validatedFiles
	// List of validated schemas
	validatedSchemas schemas
	// Init will only perform prerequisites related to elasticsearch / kafka / cockroachdb
	Init bool
	// Logger expose zerolog so it can be override
	Logger *zerolog.Logger
}

// List of validated files with SQL queries
type validatedFiles struct {
	File    string
	Queries []string
}

// createTopic reference all the possible requirements to create a topic
type createTopicModel struct {
	// Topic name
	Name string `json:"name" yaml:"name" validate:"required"`
	// Number of partitions
	NumPartitions int `json:"numPartitions" yaml:"numPartitions" validate:"required"`
	// Cluster replication factor
	ReplicationFactor int `json:"replicationFactor" yaml:"replicationFactor" validate:"required"`
	// Topic config
	TopicConfig []topicConfig
}

// topicConfig store key/value pair use to configure the topic
type topicConfig struct {
	Key   string `json:"key" yaml:"key" validate:"required"`     // Config key
	Value string `json:"value" yaml:"value" validate:"required"` // Config value
}

// kafkaWriteMessage define the requirements to write messages into kafka
type kafkaWriteMessage struct {
	TopicName string `json:"name" yaml:"name"`   // Topic name
	Key       string `json:"key" yaml:"key"`     // Message key
	Value     string `json:"value" yaml:"value"` // Message value
}

// schemas represent the global config to push data to elasticsearch
type schemas struct {
	// Config schema
	Schemas []configSchema `json:"schemas" yaml:"schemas" validate:"required,dive"`
}

// configSchema is the validator
type configSchema struct {
	// Schema name
	Name string `json:"name" yaml:"name" validate:"required"`
	// Topic schema
	Topic topicSchema `json:"topic" yaml:"topic" validate:"required"`
	// SQL requirements when advanced query is needed
	SQL sql `json:"sql" yaml:"sql" validate:"-"`
	// Elasticsearch configuration
	Elasticsearch elasticsearchSchema `json:"elasticsearch" yaml:"elasticsearch" validate:"required"`
	// // Requirements to create change feed
	ChangeFeed changeFeed `json:"changeFeed" yaml:"changeFeed" validate:"required"`
}

// topicSchema is the requirement to create the topic
type topicSchema struct {
	// Topic name
	Name string `json:"name" yaml:"name" validate:"required"`
	// Number of partitions
	NumPartitions int `json:"numPartitions" yaml:"numPartitions" validate:"required"`
	// Cluster replication factor
	ReplicationFactor int `json:"replicationFactor" yaml:"replicationFactor" validate:"required"`
	// Topic config
	TopicConfig []topicConfig `json:"config" yaml:"config" validate:"dive"`
}

// sql is the requirement to query the SQL database
type sql struct {
	// ColumnNames is the list of columns name to use to filter kafka message
	// in order to insert or delete data in elasticsearch
	Columns []string `json:"columnNames" yaml:"columnNames" validate:"required"`
	// Query is the SQL query to execute
	Query string `json:"query" yaml:"query" validate:"required"`
}

// changeFeed bind all requrirements to create changefeed
type changeFeed struct {
	// FullTableName is cockroach full table name like movr.public.promo_codes
	FullTableName string `json:"fullTableName" yaml:"fullTableName" validate:"required"`
	// Options is all change feed options required to create the change feed
	Options []string `json:"options" yaml:"options" validate:"required"`
}

// elasticsearchSchema is the requirement related to elasticsearch
type elasticsearchSchema struct {
	// Elasticsearch index
	Index elasticsearchIndex `json:"index" yaml:"index" validate:"required"`
	// Type defined if the sql query is plain or not
	Mapping map[string]interface{} `json:"mapping" yaml:"mapping" validate:"required"`
}

// elasticsearchIndex is the requirement to the elasticsearch index
type elasticsearchIndex struct {
	// Create index
	Create bool `json:"create" yaml:"create"`
	// Index name
	Name string `json:"name" yaml:"name" validate:"required"`
	// Alias name
	Alias string `json:"alias" yaml:"alias"`
}

// consumeMessage permit to consume kafka messages
type consumeMessage struct {
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
