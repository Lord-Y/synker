// Package commons assemble all functions used in other packages
package commons

import (
	"os"
	"strings"
)

// GetPGURI permit to retrieve string url to connect to the sql instance
func GetPGURI() string {
	return strings.TrimSpace(os.Getenv("SYNKER_PG_URI"))
}

// GetElasticsearchURI permit to retrieve OS env variable
func GetElasticsearchURI() string {
	return strings.TrimSpace(os.Getenv("SYNKER_ELASTICSEARCH_URI"))
}

// GetKafkaURI permit to retrieve OS env variable
func GetKafkaURI() string {
	return strings.TrimSpace(os.Getenv("SYNKER_KAFKA_URI"))
}

// GetKafkaScram permit to retrieve OS env variable
func GetKafkaScram() string {
	return strings.TrimSpace(os.Getenv("SYNKER_KAFKA_SCRAM"))
}

// GetKafkaUser permit to retrieve OS env variable
func GetKafkaUser() string {
	return strings.TrimSpace(os.Getenv("SYNKER_KAFKA_USER"))
}

// GetKafkaPassword permit to retrieve OS env variable
func GetKafkaPassword() string {
	return strings.TrimSpace(os.Getenv("SYNKER_KAFKA_PASSWORD"))
}

// GetKafkaCACert permit to retrieve OS env variable
func GetKafkaCACert() string {
	return strings.TrimSpace(os.Getenv("SYNKER_KAFKA_CACERT"))
}

// GetKafkaCert permit to retrieve OS env variable
func GetKafkaCert() string {
	return strings.TrimSpace(os.Getenv("SYNKER_KAFKA_CERT"))
}

// GetKafkaKey permit to retrieve OS env variable
func GetKafkaKey() string {
	return strings.TrimSpace(os.Getenv("SYNKER_KAFKA_KEY"))
}
