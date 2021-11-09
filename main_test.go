package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	cli "github.com/urfave/cli/v2"
)

func TestMain(t *testing.T) {
	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)

	go func() {
		<-sigc
		os.Args = []string{
			"synker",
			"-h",
		}
		main()
		signal.Stop(sigc)
	}()

	err = proc.Signal(os.Interrupt)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
}

func TestMain_fail(t *testing.T) {
	assert := assert.New(t)
	if os.Getenv("FATAL") == "1" {
		os.Args = []string{
			"synker",
			"bad",
		}
		main()
		return
	}
	cmd := exec.Command(
		os.Args[0],
		"synker",
		"bad",
		"-test.run=TestMain_fail",
	)
	cmd.Env = append(os.Environ(), "FATAL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	assert.Error(err)
}

func TestVersionDetails(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		cliArgs []string
		fail    bool
	}{
		{
			cliArgs: []string{
				"synker",
				"version-details",
			},
			fail: false,
		},
	}

	app := cli.NewApp()
	for _, tc := range tests {
		app.Writer = ioutil.Discard
		app.Commands = []*cli.Command{
			VersionDetails,
		}
		err := app.Run(tc.cliArgs)
		if tc.fail {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
	}
}

func TestValidate(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		cliArgs []string
		fail    bool
	}{
		{
			cliArgs: []string{
				"synker",
				"validate",
				"-c",
				"/tmp/test.yaml",
			},
			fail: false,
		},
	}

	app := cli.NewApp()
	for _, tc := range tests {
		app.Writer = ioutil.Discard
		app.Commands = []*cli.Command{
			CmdValidate,
		}
		err := app.Run(tc.cliArgs)
		if tc.fail {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
	}
}

func TestValidate_fail(t *testing.T) {
	assert := assert.New(t)
	if os.Getenv("FATAL") == "1" {
		os.Args = []string{
			"synker",
			"validate",
		}
		main()
		return
	}
	cmd := exec.Command(
		os.Args[0],
		"synker",
		"validate",
		"-test.run=TestValidate_fail",
	)
	cmd.Env = append(os.Environ(), "FATAL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	assert.Error(err)
}
