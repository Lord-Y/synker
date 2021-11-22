// Package validate permit to validate files provided
package validate

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	var c Validate
	c.ConfigDir = "examples/schemas"
	c.Run()
}

func TestValidate_fail_empty_dir(t *testing.T) {
	assert := assert.New(t)
	if os.Getenv("FATAL") == "1" {
		os.Args = []string{
			"synker",
			"validate",
		}
		var c Validate
		c.ConfigDir = "examples/emptydir"
		c.Run()
		return
	}
	cmd := exec.Command(
		os.Args[0],
		"synker",
		"validate",
		"-test.run=TestValidate_fail_empty_dir",
	)
	cmd.Env = append(os.Environ(), "FATAL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	assert.Error(err)
}

func TestValidate_fail_fake_dir(t *testing.T) {
	assert := assert.New(t)
	if os.Getenv("FATAL") == "1" {
		os.Args = []string{
			"synker",
			"validate",
		}
		var c Validate
		c.ConfigDir = "examples/emptydirr"
		c.Run()
		return
	}
	cmd := exec.Command(
		os.Args[0],
		"synker",
		"validate",
		"-test.run=TestValidate_fail_fake_dir",
	)
	cmd.Env = append(os.Environ(), "FATAL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	assert.Error(err)
}

func TestValidate_fail_tmp_dir(t *testing.T) {
	assert := assert.New(t)
	if os.Getenv("FATAL") == "1" {
		os.Args = []string{
			"synker",
			"validate",
		}
		var c Validate
		c.ConfigDir = "/tmp"
		c.Run()
		return
	}
	cmd := exec.Command(
		os.Args[0],
		"synker",
		"validate",
		"-test.run=TestValidate_fail_tmp_dir",
	)
	cmd.Env = append(os.Environ(), "FATAL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	assert.Error(err)
}

func TestValidate_fail_falseconfig(t *testing.T) {
	assert := assert.New(t)
	if os.Getenv("FATAL") == "1" {
		os.Args = []string{
			"synker",
			"validate",
		}
		var c Validate
		c.ConfigDir = "examples/falseconfig"
		c.Run()
		return
	}
	cmd := exec.Command(
		os.Args[0],
		"synker",
		"validate",
		"-test.run=TestValidate_fail_falseconfig",
	)
	cmd.Env = append(os.Environ(), "FATAL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	assert.Error(err)
}
