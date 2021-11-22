// Package commons assemble all functions used in other packages
package commons

import (
	"os"
	"strings"
)

// BuildDSN permit to retrieve string url to connect to the sql instance
func BuildDSN() string {
	return strings.TrimSpace(os.Getenv("SKR_PG_URI"))
}

// GetElasticsearchURI permit to retrieve OS env variable
func GetElasticsearchURI() string {
	return strings.TrimSpace(os.Getenv("SKR_ELASTICSEARCH_URI"))
}

// IsElasticsearchAuthEnabled permit to check if elasticsearch auth is enabled
func IsElasticsearchAuthEnabled() bool {
	if GetElasticsearchUser() != "" && GetElasticsearchPassword() != "" {
		return true
	}
	return false
}

// GetElasticsearchURI permit to retrieve OS env variable
func GetElasticsearchUser() string {
	return strings.TrimSpace(os.Getenv("SKR_ELASTICSEARCH_USER"))
}

// GetElasticsearchURI permit to retrieve OS env variable
func GetElasticsearchPassword() string {
	return strings.TrimSpace(os.Getenv("SKR_ELASTICSEARCH_PASSWORD"))
}

// GetKafkaURI permit to retrieve OS env variable
func GetKafkaURI() string {
	return strings.TrimSpace(os.Getenv("SKR_KAFKA_URI"))
}

// GetKafkaScram permit to retrieve OS env variable
func GetKafkaScram() string {
	return strings.TrimSpace(os.Getenv("SKR_KAFKA_SCRAM"))
}

// GetKafkaUser permit to retrieve OS env variable
func GetKafkaUser() string {
	return strings.TrimSpace(os.Getenv("SKR_KAFKA_USER"))
}

// GetKafkaPassword permit to retrieve OS env variable
func GetKafkaPassword() string {
	return strings.TrimSpace(os.Getenv("SKR_KAFKA_PASSWORD"))
}

// GetKafkaCACert permit to retrieve OS env variable
func GetKafkaCACert() string {
	return strings.TrimSpace(os.Getenv("SKR_KAFKA_CACERT"))
}

// GetKafkaCert permit to retrieve OS env variable
func GetKafkaCert() string {
	return strings.TrimSpace(os.Getenv("SKR_KAFKA_CERT"))
}

// GetKafkaKey permit to retrieve OS env variable
func GetKafkaKey() string {
	return strings.TrimSpace(os.Getenv("SKR_KAFKA_KEY"))
}
