package processing

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncreaseMetrics(t *testing.T) {
	assert := assert.New(t)

	var (
		c   Validate
		err error
	)

	os.Setenv("SYNKER_PROMETHEUS", "1")
	defer os.Unsetenv("SYNKER_PROMETHEUS")

	c.metrics, err = newMetrics()
	assert.Nil(err)

	c.increaseMetrics("kafka", "x", "y")
	c.increaseMetrics("elasticsearch", "x", "y")
}
