// Package etcd assemble all requirements to communicate with etcd cluster
package etcd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/jackc/fake"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	assert := assert.New(t)

	_, err := Client()
	assert.NoError(err)
}

func TestClient_bad_password(t *testing.T) {
	assert := assert.New(t)

	os.Setenv("SKR_ETCD_PASSWORD", "plop")
	defer os.Unsetenv("SKR_ETCD_PASSWORD")
	_, err := Client()
	assert.Error(err)
}

func TestClient_fake_ssl(t *testing.T) {
	assert := assert.New(t)
	ca, err := os.CreateTemp(os.TempDir(), fake.CharactersN(5))
	if err != nil {
		assert.Fail("Fail to create temp file")
		return
	}
	defer os.Remove(ca.Name())

	cert, err := os.CreateTemp(os.TempDir(), fake.CharactersN(5))
	if err != nil {
		assert.Fail("Fail to create temp file")
		return
	}
	defer os.Remove(cert.Name())

	key, err := os.CreateTemp(os.TempDir(), fake.CharactersN(5))
	if err != nil {
		assert.Fail("Fail to create temp file")
		return
	}
	defer os.Remove(key.Name())

	os.Setenv("SKR_ETCD_PASSWORD", "plop")
	defer os.Unsetenv("SKR_ETCD_PASSWORD")
	os.Setenv("SKR_ETCD_CACERT", ca.Name())
	os.Setenv("SKR_ETCD_CERT", cert.Name())
	os.Setenv("SKR_ETCD_KEY", key.Name())
	defer os.Unsetenv("SKR_ETCD_CACERT")
	defer os.Unsetenv("SKR_ETCD_CERT")
	defer os.Unsetenv("SKR_ETCD_KEY")
	_, err = Client()
	log.Error().Err(err).Msg("client fail")
	assert.Error(err)
}

func TestClient_bad_ssl(t *testing.T) {
	assert := assert.New(t)
	ca, err := os.CreateTemp(os.TempDir(), fake.CharactersN(5))
	if err != nil {
		assert.Fail("Fail to create temp file")
		return
	}
	defer os.Remove(ca.Name())

	cert, err := os.CreateTemp(os.TempDir(), fake.CharactersN(5))
	if err != nil {
		assert.Fail("Fail to create temp file")
		return
	}
	defer os.Remove(cert.Name())

	key, err := os.CreateTemp(os.TempDir(), fake.CharactersN(5))
	if err != nil {
		assert.Fail("Fail to create temp file")
		return
	}
	defer os.Remove(key.Name())

	cab, certb, keyb, err := certsetup()
	if err != nil {
		assert.Fail("Fail to create certificate requirements")
		return
	}
	err = ioutil.WriteFile(ca.Name(), cab.Bytes(), 0600)
	if err != nil {
		assert.Fail("Fail to write content to ca")
		return
	}
	err = ioutil.WriteFile(cert.Name(), certb.Bytes(), 0600)
	if err != nil {
		assert.Fail("Fail to write content to cert")
		return
	}
	err = ioutil.WriteFile(key.Name(), keyb.Bytes(), 0600)
	if err != nil {
		assert.Fail("Fail to write content to key")
		return
	}

	os.Setenv("SKR_ETCD_PASSWORD", "plop")
	defer os.Unsetenv("SKR_ETCD_PASSWORD")
	os.Setenv("SKR_ETCD_CACERT", ca.Name())
	os.Setenv("SKR_ETCD_CERT", cert.Name())
	os.Setenv("SKR_ETCD_KEY", key.Name())
	defer os.Unsetenv("SKR_ETCD_CACERT")
	defer os.Unsetenv("SKR_ETCD_CERT")
	defer os.Unsetenv("SKR_ETCD_KEY")
	_, err = Client()
	log.Error().Err(err).Msg("client fail")
	assert.Error(err)
}

func TestPut(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		key   string
		value string
		fail  bool
	}{
		{
			key:   "key",
			value: "test",
			fail:  false,
		},
		{
			key:   "",
			value: "value",
			fail:  true,
		},
	}

	os.Setenv("SKR_ETCD_PASSWORD", "synker")
	defer os.Unsetenv("SKR_ETCD_PASSWORD")
	cli, err := Client()
	if err != nil {
		assert.FailNow("Failed to connect to etcd cluster")
		return
	}
	for _, tc := range tests {
		_, err := Put(cli, tc.key, tc.value)
		if tc.fail {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
	}
}

func TestGet(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		key  string
		fail bool
	}{
		{
			key:  "key",
			fail: false,
		},
		{
			key:  "",
			fail: true,
		},
	}

	os.Setenv("SKR_ETCD_PASSWORD", "synker")
	defer os.Unsetenv("SKR_ETCD_PASSWORD")
	cli, err := Client()
	if err != nil {
		assert.FailNow("Failed to connect to etcd cluster")
		return
	}
	for _, tc := range tests {
		_, err := Get(cli, tc.key)
		if tc.fail {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
	}
}

func TestWatch(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		key   string
		value string
		fail  bool
	}{
		{
			key:   "key",
			value: "test",
			fail:  false,
		},
		{
			key:   "",
			value: "value",
			fail:  true,
		},
	}

	os.Setenv("SKR_ETCD_PASSWORD", "synker")
	defer os.Unsetenv("SKR_ETCD_PASSWORD")
	os.Setenv("SKR_ETCD_PASSWORD", "synker")
	defer os.Unsetenv("SKR_ETCD_PASSWORD")
	cli, err := Client()
	if err != nil {
		assert.FailNow("Failed to connect to etcd cluster")
		return
	}
	defer cli.Close()
	cliput, err := Client()
	if err != nil {
		assert.FailNow("Failed to connect to etcd cluster")
		return
	}
	for _, tc := range tests {
		if tc.fail {
			_, err := Watch(cli, tc.key)
			assert.Error(err)
		} else {
			resp, err := Watch(cli, tc.key)
			assert.NoError(err)
			_, err = Put(cliput, tc.key, tc.value)
			assert.NoError(err)
			assert.Len(resp, 0)
		}
	}
}

func TestWatchWithPrefix(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		key   string
		value string
		fail  bool
	}{
		{
			key:   "key",
			value: "test",
			fail:  false,
		},
		{
			key:   "",
			value: "value",
			fail:  true,
		},
	}

	os.Setenv("SKR_ETCD_PASSWORD", "synker")
	defer os.Unsetenv("SKR_ETCD_PASSWORD")
	os.Setenv("SKR_ETCD_PASSWORD", "synker")
	defer os.Unsetenv("SKR_ETCD_PASSWORD")
	cli, err := Client()
	if err != nil {
		assert.FailNow("Failed to connect to etcd cluster")
		return
	}
	defer cli.Close()
	cliput, err := Client()
	if err != nil {
		assert.FailNow("Failed to connect to etcd cluster")
		return
	}
	for _, tc := range tests {
		if tc.fail {
			_, err := WatchWithPrefix(cli, tc.key)
			assert.Error(err)
		} else {
			resp, err := WatchWithPrefix(cli, tc.key)
			assert.NoError(err)
			_, err = Put(cliput, tc.key, tc.value)
			assert.NoError(err)
			assert.Len(resp, 0)
		}
	}
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		key  string
		fail bool
	}{
		{
			key:  "key",
			fail: false,
		},
		{
			key:  "",
			fail: true,
		},
	}

	os.Setenv("SKR_ETCD_PASSWORD", "synker")
	defer os.Unsetenv("SKR_ETCD_PASSWORD")
	cli, err := Client()
	if err != nil {
		assert.FailNow("Failed to connect to etcd cluster")
		return
	}
	for _, tc := range tests {
		_, err := Delete(cli, tc.key)
		if tc.fail {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
	}
}
