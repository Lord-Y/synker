// Package validate permit to validate files provided
package validate

import (
	"os"
	"os/exec"
	"os/signal"
	"testing"
	"time"

	"github.com/Lord-Y/synker/elasticsearch"
	"github.com/Lord-Y/synker/kafka"
	"github.com/Lord-Y/synker/tools"
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

func TestManageTopicsAndElasticsearchIndex(t *testing.T) {
	assert := assert.New(t)

	conn, err := kafka.Client()
	assert.Nil(err)
	topics, err := kafka.ListTopics(conn)
	assert.Nil(err)
	if tools.InSlice("movr.public.user_promo_codes", topics) {
		conndel, err := kafka.Client()
		assert.Nil(err)
		err = kafka.DeleteTopics(
			conndel,
			[]string{
				"movr.public.user_promo_codes",
			},
		)
		assert.Nil(err)
		client, err := elasticsearch.Client()
		assert.Nil(err)

		b, err := elasticsearch.DeleteIndex(client, "user_promo_codes")
		assert.Nil(err)
		assert.Equal(true, b)
	}

	var c Validate
	c.ConfigDir = "examples/schemas"
	c.Run()

	err = c.ManageTopics()
	assert.Nil(err)
	err = c.ManageElasticsearchIndex()
	assert.Nil(err)
}

func TestManageTopicsAndElasticsearchIndex_with_alias(t *testing.T) {
	assert := assert.New(t)

	conn, err := kafka.Client()
	assert.Nil(err)
	topics, err := kafka.ListTopics(conn)
	assert.Nil(err)
	if tools.InSlice("movr.public.user_promo_codes", topics) {
		conndel, err := kafka.Client()
		assert.Nil(err)
		err = kafka.DeleteTopics(
			conndel,
			[]string{
				"movr.public.user_promo_codes",
			},
		)
		assert.Nil(err)
		client, err := elasticsearch.Client()
		assert.Nil(err)

		b, err := elasticsearch.DeleteIndex(client, "user_promo_codes")
		assert.Nil(err)
		assert.Equal(true, b)
	}

	var c Validate
	c.ConfigDir = "examples/schemas"
	c.Run()

	err = c.ManageTopics()
	assert.Nil(err)

	var schema_id int
	for k := range c.ValidatedSchemas.Schemas {
		if c.ValidatedSchemas.Schemas[k].Name == "user_promo_codes" {
			schema_id = k
			break
		}
	}
	c.ValidatedSchemas.Schemas[schema_id].Elasticsearch.Index.Alias = "user_promo_codes_alias"
	err = c.ManageElasticsearchIndex()
	assert.Nil(err)
}

func TestManageTopicsAndElasticsearchIndex_with_same_index_alias(t *testing.T) {
	assert := assert.New(t)

	conn, err := kafka.Client()
	assert.Nil(err)
	topics, err := kafka.ListTopics(conn)
	assert.Nil(err)
	if tools.InSlice("movr.public.user_promo_codes", topics) {
		conndel, err := kafka.Client()
		assert.Nil(err)
		err = kafka.DeleteTopics(
			conndel,
			[]string{
				"movr.public.user_promo_codes",
			},
		)
		assert.Nil(err)
		client, err := elasticsearch.Client()
		assert.Nil(err)

		b, err := elasticsearch.DeleteIndex(client, "user_promo_codes")
		assert.Nil(err)
		assert.Equal(true, b)
	}

	var c Validate
	c.ConfigDir = "examples/schemas"
	c.Run()

	err = c.ManageTopics()
	assert.Nil(err)

	var schema_id int
	for k := range c.ValidatedSchemas.Schemas {
		if c.ValidatedSchemas.Schemas[k].Name == "user_promo_codes" {
			schema_id = k
			break
		}
	}
	c.ValidatedSchemas.Schemas[schema_id].Elasticsearch.Index.Alias = "user_promo_codes"
	err = c.ManageElasticsearchIndex()
	assert.Error(err)
}

func TestManageTopicsAndElasticsearchIndex_with_same_index_alias_after_index_created(t *testing.T) {
	assert := assert.New(t)

	conn, err := kafka.Client()
	assert.Nil(err)
	topics, err := kafka.ListTopics(conn)
	assert.Nil(err)
	if tools.InSlice("movr.public.user_promo_codes", topics) {
		conndel, err := kafka.Client()
		assert.Nil(err)
		err = kafka.DeleteTopics(
			conndel,
			[]string{
				"movr.public.user_promo_codes",
			},
		)
		assert.Nil(err)
		client, err := elasticsearch.Client()
		assert.Nil(err)

		b, err := elasticsearch.DeleteIndex(client, "user_promo_codes")
		assert.Nil(err)
		assert.Equal(true, b)
	}

	var c Validate
	c.ConfigDir = "examples/schemas"
	c.Run()

	err = c.ManageTopics()
	assert.Nil(err)

	var schema_id int
	for k := range c.ValidatedSchemas.Schemas {
		if c.ValidatedSchemas.Schemas[k].Name == "user_promo_codes" {
			schema_id = k
			break
		}
	}
	err = c.ManageElasticsearchIndex()
	assert.Nil(err)
	c.ValidatedSchemas.Schemas[schema_id].Elasticsearch.Index.Alias = "user_promo_codes"
	err = c.ManageElasticsearchIndex()
	assert.Error(err)
}

func TestValidate_processing(t *testing.T) {
	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)

	go func() {
		<-sigc
		var c Validate
		c.ConfigDir = "examples/schemas"
		c.Run()
		c.Processing()
		signal.Stop(sigc)
	}()

	err = proc.Signal(os.Interrupt)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
}
