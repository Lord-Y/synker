// Package processing provide all requirements to process change data capture
package processing

import (
	"context"
	"os"
	"os/exec"
	"os/signal"
	"testing"
	"time"

	"github.com/Lord-Y/synker/commons"
	"github.com/Lord-Y/synker/elasticsearch"
	"github.com/Lord-Y/synker/kafka"
	"github.com/Lord-Y/synker/tools"
	"github.com/icrowley/fake"
	"github.com/jackc/pgx/v5/pgxpool"
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

		indexExist, err := elasticsearch.IndexAlreadyExist(client, "user_promo_codes")
		assert.Nil(err)
		if indexExist {
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

func TestManageChangeFeed(t *testing.T) {
	assert := assert.New(t)

	var c Validate
	c.ConfigDir = "examples/schemas"
	c.Run()
	err := c.ManageChangeFeed()
	assert.Nil(err)
}

func TestValidate_processing(t *testing.T) {
	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)

	go func() {
		var c Validate
		c.ConfigDir = "examples/schemas"
		c.Run()
		go func() {
			time.Sleep(60 * time.Second)
			<-sigc
			signal.Stop(sigc)
		}()
		c.Processing()
	}()

	err = proc.Signal(os.Interrupt)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
}

func (c *Validate) testProcessing(stop chan struct{}) {
	var i int
	time.AfterFunc(90*time.Second, func() {
		stop <- struct{}{}
	})

	for {
		select {
		case <-stop:
			return
		default:
			if i == 0 {
				i++
				c.Processing()
			}
		}
	}
}

func TestValidate_processing_with_new_id(t *testing.T) {
	assert := assert.New(t)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)

	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}

	fakeCharacter := fake.CharactersN(5)
	var c Validate
	c.ConfigDir = "examples/schemas"
	c.Run()
	err = manageUserPromoCodesForUnitTesting("add", fakeCharacter)
	assert.Nil(err)

	stop := make(chan struct{})
	go c.testProcessing(stop)

	time.Sleep(40 * time.Second)
	err = manageUserPromoCodesForUnitTesting("delete", fakeCharacter)
	assert.Nil(err)
	<-stop

	if err := proc.Signal(os.Interrupt); err != nil {
		t.Fatal(err)
	}
}

func manageUserPromoCodesForUnitTesting(action, fakeCharacter string) (err error) {
	cfg, err := pgxpool.ParseConfig(commons.GetPGURI())
	if err != nil {
		return
	}

	ctx := context.Background()
	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return
	}
	defer db.Close()

	tx, err := db.Begin(ctx)
	if err != nil {
		return
	}
	//golangci-lint fail on this check while the transaction error is checked
	defer tx.Rollback(ctx) //nolint

	switch action {
	case "delete":
		_, err = tx.Exec(
			ctx,
			"DELETE FROM user_promo_codes WHERE city = $1 AND user_id = $2 AND code = $3 AND usage_count = $4",
			"amsterdam",
			"ae147ae1-47ae-4800-8000-000000000022",
			fakeCharacter,
			10,
		)
	default:
		_, err = tx.Exec(
			ctx,
			"INSERT INTO user_promo_codes VALUES($1,$2,$3,NOW(),$4)",
			"amsterdam",
			"ae147ae1-47ae-4800-8000-000000000022",
			fakeCharacter,
			10,
		)
	}
	if err != nil {
		return
	}

	if err = tx.Commit(ctx); err != nil {
		return
	}
	return
}
