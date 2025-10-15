package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
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
	if os.Getenv("FATAL") == "1" {
		main()
		return
	}

	pwd, _ := os.Getwd()
	cmd := exec.Command(
		filepath.Join(pwd, os.Args[0]),
		"bad",
		"-test.run=TestMain_fail",
	)
	cmd.Env = append(os.Environ(), "FATAL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, wanted exit status 1", err)
}

func TestMain_fail_api(t *testing.T) {
	if os.Getenv("FATAL") == "1" {
		main()
		return
	}

	pwd, _ := os.Getwd()
	cmd := exec.Command(
		filepath.Join(pwd, os.Args[0]),
		"api",
		"-c",
		"-test.run=TestMain_fail_api",
	)
	cmd.Env = append(os.Environ(), "FATAL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, wanted exit status 1", err)
}

func TestCommons(t *testing.T) {
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
		{
			cliArgs: []string{
				"synker",
				"validate",
				"-c",
				"processing/examples/schemas",
			},
			fail: false,
		},
		{
			cliArgs: []string{
				"synker",
				"validate",
			},
			fail: true,
		},
	}

	for _, tc := range tests {
		cmd := cli.Command{}
		cmd.Writer = io.Discard
		cmd.Commands = []*cli.Command{
			VersionDetails,
			CmdValidate,
		}
		err := cmd.Run(context.Background(), tc.cliArgs)
		if tc.fail {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
	}
}

type os_variable struct {
	name, option              string
	t                         *testing.T
	vars, varsBackup, cliArgs []string
	k                         int
}

func (test *os_variable) testMain_fail_os_variables() {
	test.t.Run(fmt.Sprintf("%s_%s_%s", test.name, test.option, test.vars[test.k]), func(t *testing.T) {
		if os.Getenv(fmt.Sprintf("FATAL_%s_%s", test.option, test.vars[test.k])) == "1" {
			os.Args = test.cliArgs
			main()
			return
		}

		cmd := exec.Command(
			os.Args[0],
			"-test.run="+test.name,
		)
		cmd.Env = append(os.Environ(), fmt.Sprintf("FATAL_%s_%s=1", test.option, test.vars[test.k]))
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			return
		}
		t.Fatalf("process ran with err %v, wanted exit status 1", err)
	})
}

func TestMain_fail_os_variables(t *testing.T) {
	var os_variable os_variable

	os_variable.t = t
	os_variable.name = "TestMain_fail_os_variables"
	os_variable.vars = []string{
		"SYNKER_PG_URI",
		"SYNKER_ELASTICSEARCH_URI",
		"SYNKER_KAFKA_URI",
	}

	for _, os_var := range os_variable.vars {
		os_variable.varsBackup = append(os_variable.varsBackup, os.Getenv(os_var))
	}

	for _, option := range []string{"api", "init"} {
		for k := range os_variable.vars {
			cliArgs := []string{
				"",
				option,
				"-c",
				"processing/examples/schemas",
			}
			os_variable.k = k
			os_variable.cliArgs = cliArgs
			os_variable.option = option
			_ = os.Unsetenv(os_variable.vars[k])
			if k > 0 {
				for index := range len(os_variable.vars) {
					if k != index {
						_ = os.Setenv(os_variable.vars[index], os_variable.varsBackup[index])
					}
				}
			}
			os_variable.testMain_fail_os_variables()
		}
	}
}
