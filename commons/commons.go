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
