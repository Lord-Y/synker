// Package processing provide all requirements to process change data capture
package processing

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/Lord-Y/synker/logger"
	"github.com/Lord-Y/synker/tools"
	"github.com/stretchr/testify/assert"
)

func TestParseAndValidateConfig(t *testing.T) {
	assert := assert.New(t)
	var c Validate
	c.Logger = logger.NewLogger()
	c.ConfigDir = "examples/schemas"
	assert.Nil(c.ParseAndValidateConfig())
}

func TestParsing_failures(t *testing.T) {
	assert := assert.New(t)

	funcName := "TestParsing_failures"
	dirs := []string{
		"examples/emptydir",
		"examples/emptydirr",
		"/tmp",
		"examples/falseconfig",
		"examples/doesnotexist",
	}

	for i, dir := range dirs {
		if os.Getenv(fmt.Sprintf("%s_%d", funcName, i)) == "1" {
			var c Validate
			c.Logger = logger.NewLogger()
			c.ConfigDir = dir
			assert.Nil(c.ParseAndValidateConfig())
			return
		}
		cmd := exec.Command(
			os.Args[0],
			"-test.run="+funcName,
		)
		cmd.Env = append(os.Environ(), fmt.Sprintf("%s_%d", funcName, i)+"=1")
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			return
		}
		t.Fatalf("process ran with err %v, wanted exit status 1", err)
	}
}

func TestParsing_falseconfig(t *testing.T) {
	assert := assert.New(t)

	var (
		c     Validate
		files []parsedFile
	)

	dirs := []string{
		"examples/emptydir",
		"examples/falseconfig",
	}

	c.Logger = logger.NewLogger()
	for _, dir := range dirs {
		c.ConfigDir = dir
		err := filepath.WalkDir(
			dir,
			func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if tools.InSlice(filepath.Ext(d.Name()), valideFiles) {
					files = append(files, parsedFile{
						file:   path,
						isJson: filepath.Ext(d.Name()) == ".json",
					})
				}
				return nil
			},
		)
		assert.Nil(err)

		for _, v := range files {
			c.files = append(c.files, v)
			_, err = c.parsing()
			assert.Error(err)
			c.files = nil
		}
	}
}
