// Package tls assemble all requirements to create tls ca, cert and key
package tls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCertSetup(t *testing.T) {
	assert := assert.New(t)

	_, _, _, err := CertSetup()
	assert.Nil(err)
}
