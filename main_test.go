package main

import (
	"fmt"
	"io"
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

func TestMain_init(t *testing.T) {
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
			"init",
			"-c",
			"processing/examples/schemas",
		}
		main()
		signal.Stop(sigc)
	}()
	time.Sleep(1 * time.Minute)

	err = proc.Signal(os.Interrupt)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
}

func TestMain_api(t *testing.T) {
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
			"api",
			"-c",
			"processing/examples/schemas",
		}
		main()
		signal.Stop(sigc)
	}()
	time.Sleep(1 * time.Minute)

	err = proc.Signal(os.Interrupt)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
}

func TestMain_api_init(t *testing.T) {
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
			"api",
			"-c",
			"processing/examples/schemas",
			"-i",
		}
		main()
		signal.Stop(sigc)
	}()
	time.Sleep(1 * time.Minute)

	err = proc.Signal(os.Interrupt)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
}

func TestMain_api_port(t *testing.T) {
	ports := []string{"18080", ":18080"}
	for k := range ports {
		os.Setenv("SYNKER_API_PORT", ports[k])
		defer os.Unsetenv("SYNKER_API_PORT")

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
				"api",
				"-c",
				"processing/examples/schemas",
			}
			main()
			signal.Stop(sigc)
		}()
		time.Sleep(1 * time.Minute)

		err = proc.Signal(os.Interrupt)
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(1 * time.Second)
	}
}

func TestMain_fail(t *testing.T) {
	assert := assert.New(t)
	if os.Getenv("FATAL") == "1" {
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

func TestMain_fail_api(t *testing.T) {
	assert := assert.New(t)
	if os.Getenv("FATAL") == "1" {
		main()
		return
	}
	cmd := exec.Command(
		os.Args[0],
		"synker",
		"api",
		"-c",
		"-test.run=TestMain_fail_api",
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
		app.Writer = io.Discard
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
				"processing/examples/schemas",
			},
			fail: false,
		},
	}

	app := cli.NewApp()
	for _, tc := range tests {
		app.Writer = io.Discard
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

func TestValidate_fail_empty_dir(t *testing.T) {
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
		"-test.run=TestValidate_fail_empty_dir",
	)
	cmd.Env = append(os.Environ(), "FATAL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	assert.Error(err)
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

type api_os_variable struct {
	t    *testing.T
	vars []string
	k    int
}

func (test *api_os_variable) testMain_fail_api_os_variables() {
	test.t.Run(fmt.Sprintf("testMain_fail_api_os_variables_%s", test.vars[test.k]), func(t *testing.T) {
		assert := assert.New(t)
		if os.Getenv("FATAL") == "1" {
			main()
			return
		}

		cmd := exec.Command(
			os.Args[0],
			"synker",
			"api",
			"-c",
			"processing/examples/schemas",
			"-test.run=TestMain_fail_api_os_variables",
		)

		switch test.k {
		case 0:
			_ = os.Unsetenv(test.vars[test.k])
		case 1:
			os.Setenv("SYNKER_PG_URI", "postgres://root:@127.0.0.1:26257/movr?sslmode=disable")
			_ = os.Unsetenv(test.vars[test.k])
		case 2:
			os.Setenv("SYNKER_ELASTICSEARCH_URI", "http://127.0.0.1:9200")
			_ = os.Unsetenv(test.vars[test.k])
		}

		cmd.Env = append(os.Environ(), "FATAL=1")
		err := cmd.Run()
		assert.Error(err)
	})
}

func TestMain_fail_api_os_variables(t *testing.T) {
	var api_os_variable api_os_variable

	api_os_variable.t = t
	api_os_variable.vars = []string{
		"SYNKER_PG_URI",
		"SYNKER_ELASTICSEARCH_URI",
		"SYNKER_KAFKA_URI",
	}

	for k := range api_os_variable.vars {
		api_os_variable.k = k
		api_os_variable.testMain_fail_api_os_variables()
	}
	os.Setenv("SYNKER_KAFKA_URI", "localhost:9092")
}
