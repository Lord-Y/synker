// Package models assemble all structs, interface e.g ...
package models

// Configuration reference all requirements relation to the file config
type Configuration struct {
	ConfigDir string `json:"configDir" yaml:"configDir"` // config dir path
}

// CreateTopic reference all the possible requirements to create a topic
type CreateTopic struct {
	Name              string `json:"name" yaml:"name"`                           // topic name
	NumPartitions     int    `json:"numPartitions" yaml:"numPartitions"`         // partition number
	ReplicationFactor int    `json:"replicationFactor" yaml:"replicationFactor"` // replication factor
	TopicConfig       []TopicConfig
}

// TopicConfig store key/value pair use to configure the topic
type TopicConfig struct {
	Key   string `json:"key" yaml:"key"`     // config key
	Value string `json:"value" yaml:"value"` // config value
}

// KafkaWriteMessage define the requirements to write messages into kafka
type KafkaWriteMessage struct {
	TopicName string `json:"name" yaml:"name"`   // topic name
	Key       string `json:"key" yaml:"key"`     // message key
	Value     string `json:"value" yaml:"value"` // message value
}
