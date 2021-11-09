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

// GetEleasticsearchURI permit to retrieve OS env variable
func GetEleasticsearchURI() string {
	return strings.TrimSpace(os.Getenv("SKR_ELASTICSEARCH_URI"))
}

// GetEtcdURI permit to retrieve OS env variable
func GetEtcdURI() string {
	return strings.TrimSpace(os.Getenv("SKR_ETCD_URI"))
}

// GetEtcdUser permit to retrieve OS env variable
func GetEtcdUser() string {
	return strings.TrimSpace(os.Getenv("SKR_ETCD_USER"))
}

// GetEtcdPassword permit to retrieve OS env variable
func GetEtcdPassword() string {
	return strings.TrimSpace(os.Getenv("SKR_ETCD_PASSWORD"))
}

// GetEtcdCACert permit to retrieve OS env variable
func GetEtcdCACert() string {
	return strings.TrimSpace(os.Getenv("SKR_ETCD_CACERT"))
}

// GetEtcdCert permit to retrieve OS env variable
func GetEtcdCert() string {
	return strings.TrimSpace(os.Getenv("SKR_ETCD_CERT"))
}

// GetEtcdKey permit to retrieve OS env variable
func GetEtcdKey() string {
	return strings.TrimSpace(os.Getenv("SKR_ETCD_KEY"))
}
