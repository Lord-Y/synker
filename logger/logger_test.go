package logger

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		logLevel string
		expected string
	}{
		{
			logLevel: "info",
			expected: "info",
		},
		{
			logLevel: "warn",
			expected: "warn",
		},
		{
			logLevel: "debug",
			expected: "debug",
		},
		{
			logLevel: "error",
			expected: "error",
		},
		{
			logLevel: "fatal",
			expected: "fatal",
		},
		{
			logLevel: "trace",
			expected: "trace",
		},
		{
			logLevel: "panic",
			expected: "panic",
		},
		{
			logLevel: "plop",
			expected: "info",
		},
	}

	for _, tc := range tests {
		os.Setenv("SYNKER_LOG_LEVEL", tc.logLevel)
		log.Logger = *NewLogger()

		assert.Equal(tc.expected, zerolog.GlobalLevel().String())
		os.Unsetenv("SYNKER_LOG_LEVEL")

		os.Setenv("SYNKER_LOG_LEVEL", tc.logLevel)
		os.Setenv("SYNKER_LOG_FORMAT_JSON", "true")
		defer os.Unsetenv("SYNKER_LOG_FORMAT_JSON")
		log.Logger = *NewLogger()

		assert.Equal(tc.expected, zerolog.GlobalLevel().String())
		os.Unsetenv("SYNKER_LOG_LEVEL")
	}
}

func TestNewLogger_info(t *testing.T) {
	log.Logger = *NewLogger()
	log.Info().Msgf("Testing logger")
}
