package processing

import (
	"testing"

	"github.com/Lord-Y/synker/logger"
	"github.com/stretchr/testify/assert"
)

func TestManageChangeFeed(t *testing.T) {
	assert := assert.New(t)

	var c Validate
	c.Logger = logger.NewLogger()
	c.ConfigDir = "examples/schemas"
	assert.Nil(c.ParseAndValidateConfig())
	assert.Nil(c.manageChangeFeed())
}
