package logger

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestSetCLILoggerLogLevel(t *testing.T) {
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
		os.Setenv("SKR_LOG_LEVEL", tc.logLevel)
		SetCLILoggerLogLevel()
		z := zerolog.GlobalLevel().String()

		assert.Equal(tc.expected, z)
		os.Unsetenv("SKR_LOG_LEVEL")
	}
}

func TestCLILogger_info(t *testing.T) {
	SetCLILoggerLogLevel()
	log.Info().Msgf("Testing logger")
}

func TestSetAPILoggerLogLevel(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		logLevel string
		expected string
		errorMsg string
	}{
		{
			logLevel: "info",
			expected: "info",
			errorMsg: "Fail to get expected logLevel",
		},
		{
			logLevel: "warn",
			expected: "warn",
			errorMsg: "Fail to get expected logLevel",
		},
		{
			logLevel: "debug",
			expected: "debug",
			errorMsg: "Fail to get expected logLevel",
		},
		{
			logLevel: "error",
			expected: "error",
			errorMsg: "Fail to get expected logLevel",
		},
		{
			logLevel: "fatal",
			expected: "fatal",
			errorMsg: "Fail to get expected logLevel",
		},
		{
			logLevel: "trace",
			expected: "trace",
			errorMsg: "Fail to get expected logLevel",
		},
		{
			logLevel: "panic",
			expected: "panic",
			errorMsg: "Fail to get expected logLevel",
		},
		{
			logLevel: "plop",
			expected: "info",
			errorMsg: "Fail to get expected logLevel",
		},
	}

	for _, tc := range tests {
		os.Setenv("SKR_LOG_LEVEL", tc.logLevel)
		SetAPILoggerLogLevel()
		z := zerolog.GlobalLevel().String()

		assert.Equal(tc.expected, z)
		os.Unsetenv("SKR_LOG_LEVEL")
	}
}
