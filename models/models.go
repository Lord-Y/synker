// Package models assemble all structs, interface e.g ...
package models

// Configuration reference all requirements relation to the file config
type Configuration struct {
	ConfigFile string `json:"config" yaml:"config"` // config file path
}
