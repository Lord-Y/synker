// Package validate permit to validate files provided
package validate

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Lord-Y/synker/logger"
	"github.com/Lord-Y/synker/models"
	"github.com/Lord-Y/synker/tools"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Validate models.Configuration

var (
	// valideFiles will be used when walking into the directory tree
	valideFiles = []string{
		".json",
		".yaml",
		".yml",
	}
)

func init() {
	logger.SetCLILoggerLogLevel()
}

// Run will validate the files configurations
func (c *Validate) Run() {
	var files []string
	err := filepath.WalkDir(
		c.ConfigDir,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if tools.InSlice(filepath.Ext(d.Name()), valideFiles) {
				files = append(files, path)
			}
			return nil
		},
	)
	if err != nil {
		log.Error().Err(err).Msgf("Fail to walk into directory %s", c.ConfigDir)
		return
	}
	if len(files) == 0 {
		log.Fatal().Msg("No config files found to validate")
		return
	}
	c.Files = files
	z, err := c.parsing()
	if err != nil {
		log.Error().Err(err).Msgf("File %s is invalid", z)
		return
	}
	for _, v := range c.ValidatedFiles {
		log.Info().Msgf("Config file %s is valid", v)
	}
}

// loadFiles will load all files requests and return a slice of bytes
func loadFiles(f string) (z []byte, err error) {
	of, err := os.Open(f)
	if err != nil {
		return z, fmt.Errorf("Fail to open file `%s` on your system", f)
	}
	defer of.Close()
	fo, err := ioutil.ReadAll(of)
	if err != nil {
		return z, fmt.Errorf("Fail to read file content `%s` on your system", f)
	}
	return fo, err
}

// parsing permit to validate all provided config
func (c *Validate) parsing() (file string, err error) {
	for _, file = range c.Files {
		fBytes, err := loadFiles(file)
		if err != nil {
			return file, err
		}
		var z models.Schemas
		if tools.IsYamlFromBytes(fBytes) {
			err = yaml.Unmarshal(fBytes, &z)
			if err != nil {
				return file, err
			}
		}
		if tools.IsJSONFromBytes(fBytes) {
			err = json.Unmarshal(fBytes, &z)
			if err != nil {
				return file, err
			}
		}
		c.ValidatedFiles = append(c.ValidatedFiles, file)
		if len(c.ValidatedSchemas.Schemas) == 0 {
			c.ValidatedSchemas = z
		} else {
			c.ValidatedSchemas.Schemas = append(c.ValidatedSchemas.Schemas, z.Schemas...)
		}
	}
	return
}
