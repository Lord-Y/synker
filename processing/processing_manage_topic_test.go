package processing

import (
	"testing"

	"github.com/Lord-Y/synker/logger"
	"github.com/Lord-Y/synker/tools"
	"github.com/stretchr/testify/assert"
)

func TestManageTopicsAndElasticsearchIndex(t *testing.T) {
	assert := assert.New(t)
	var c Validate
	c.Logger = logger.NewLogger()

	conn, err := c.kClient()
	assert.Nil(err)
	topics, err := c.listTopics(conn)
	assert.Nil(err)
	if tools.InSlice("movr.public.user_promo_codes", topics) {
		conndel, err := c.kClient()
		assert.Nil(err)
		err = c.deleteTopics(
			conndel,
			[]string{
				"movr.public.user_promo_codes",
			},
		)
		assert.Nil(err)
		client, err := c.eClient()
		assert.Nil(err)

		indexExist, err := c.indexAlreadyExist(client, "user_promo_codes")
		assert.Nil(err)
		if indexExist {
			client, err := c.eClient()
			assert.Nil(err)
			b, err := c.deleteIndex(client, "user_promo_codes")
			assert.Nil(err)
			assert.Equal(true, b)
		}

		c.ConfigDir = "examples/schemas"
		assert.Nil(c.ParseAndValidateConfig())

		err = c.manageTopics()
		assert.Nil(err)
		err = c.manageElasticsearchIndex()
		assert.Nil(err)
	}
}

func TestManageTopicsAndElasticsearchIndex_with_alias(t *testing.T) {
	assert := assert.New(t)
	var c Validate
	c.Logger = logger.NewLogger()

	conn, err := c.kClient()
	assert.Nil(err)
	topics, err := c.listTopics(conn)
	assert.Nil(err)
	if tools.InSlice("movr.public.user_promo_codes", topics) {
		conndel, err := c.kClient()
		assert.Nil(err)
		err = c.deleteTopics(
			conndel,
			[]string{
				"movr.public.user_promo_codes",
			},
		)
		assert.Nil(err)
		client, err := c.eClient()
		assert.Nil(err)

		b, err := c.deleteIndex(client, "user_promo_codes")
		assert.Nil(err)
		assert.Equal(true, b)
	}

	c.ConfigDir = "examples/schemas"
	assert.Nil(c.ParseAndValidateConfig())

	err = c.manageTopics()
	assert.Nil(err)

	var schema_id int
	for k := range c.validatedSchemas.Schemas {
		if c.validatedSchemas.Schemas[k].Name == "user_promo_codes" {
			schema_id = k
			break
		}
	}
	c.validatedSchemas.Schemas[schema_id].Elasticsearch.Index.Alias = "user_promo_codes_alias"
	err = c.manageElasticsearchIndex()
	assert.Nil(err)
}

func TestManageTopicsAndElasticsearchIndex_with_same_index_alias(t *testing.T) {
	assert := assert.New(t)
	var c Validate
	c.Logger = logger.NewLogger()

	conn, err := c.kClient()
	assert.Nil(err)
	topics, err := c.listTopics(conn)
	assert.Nil(err)
	if tools.InSlice("movr.public.user_promo_codes", topics) {
		conndel, err := c.kClient()
		assert.Nil(err)
		err = c.deleteTopics(
			conndel,
			[]string{
				"movr.public.user_promo_codes",
			},
		)
		assert.Nil(err)
		client, err := c.eClient()
		assert.Nil(err)

		b, err := c.deleteIndex(client, "user_promo_codes")
		assert.Nil(err)
		assert.Equal(true, b)
	}

	c.ConfigDir = "examples/schemas"
	assert.Nil(c.ParseAndValidateConfig())

	err = c.manageTopics()
	assert.Nil(err)

	var schema_id int
	for k := range c.validatedSchemas.Schemas {
		if c.validatedSchemas.Schemas[k].Name == "user_promo_codes" {
			schema_id = k
			break
		}
	}
	c.validatedSchemas.Schemas[schema_id].Elasticsearch.Index.Alias = "user_promo_codes"
	err = c.manageElasticsearchIndex()
	assert.Error(err)
}

func TestManageTopicsAndElasticsearchIndex_with_same_index_alias_after_index_created(t *testing.T) {
	assert := assert.New(t)
	var c Validate
	c.Logger = logger.NewLogger()

	conn, err := c.kClient()
	assert.Nil(err)
	topics, err := c.listTopics(conn)
	assert.Nil(err)
	if tools.InSlice("movr.public.user_promo_codes", topics) {
		conndel, err := c.kClient()
		assert.Nil(err)
		err = c.deleteTopics(
			conndel,
			[]string{
				"movr.public.user_promo_codes",
			},
		)
		if err != nil {
			c.Logger.Error().Err(err).Msg("TestManageTopicsAndElasticsearchIndex_with_same_index_alias_after_index_created")
		}
		assert.Nil(err)
		client, err := c.eClient()
		assert.Nil(err)

		b, err := c.deleteIndex(client, "user_promo_codes")
		assert.Nil(err)
		assert.Equal(true, b)
	}

	c.ConfigDir = "examples/schemas"
	assert.Nil(c.ParseAndValidateConfig())

	err = c.manageTopics()
	assert.Nil(err)

	var schema_id int
	for k := range c.validatedSchemas.Schemas {
		if c.validatedSchemas.Schemas[k].Name == "user_promo_codes" {
			schema_id = k
			break
		}
	}
	err = c.manageElasticsearchIndex()
	assert.Nil(err)
	c.validatedSchemas.Schemas[schema_id].Elasticsearch.Index.Alias = "user_promo_codes"
	err = c.manageElasticsearchIndex()
	assert.Error(err)
}
