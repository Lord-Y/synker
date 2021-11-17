// Package models assemble all structs, interface e.g ...
package models

// Configuration reference all requirements relation to the file config
type Configuration struct {
	ConfigDir string `json:"configDir" yaml:"configDir"` // config dir path
}
