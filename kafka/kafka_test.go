// Kafka package will handle all kafka requirements
package kafka

import (
	"os"
	"os/signal"
	"testing"
	"time"

	"github.com/Lord-Y/synker/models"
	"github.com/Lord-Y/synker/tls"
	"github.com/jackc/fake"
	"github.com/stretchr/testify/assert"
)

var (
	test_create_topic string = "test_create_topic"
)

func TestClient(t *testing.T) {
	assert := assert.New(t)

	_, err := Client()
	assert.Nil(err)
}

func TestClient_tls_bad(t *testing.T) {
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

	cab, certb, keyb, err := tls.CertSetup()
	if err != nil {
		assert.Fail("Fail to create certificate requirements")
		return
	}
	err = os.WriteFile(ca.Name(), cab.Bytes(), 0600)
	if err != nil {
		assert.Fail("Fail to write content to ca")
		return
	}
	err = os.WriteFile(cert.Name(), certb.Bytes(), 0600)
	if err != nil {
		assert.Fail("Fail to write content to cert")
		return
	}
	err = os.WriteFile(key.Name(), keyb.Bytes(), 0600)
	if err != nil {
		assert.Fail("Fail to write content to key")
		return
	}

	os.Setenv("SYNKER_KAFKA_SCRAM", "youhou")
	os.Setenv("SYNKER_KAFKA_USER", "youhou")
	os.Setenv("SYNKER_KAFKA_PASSWORD", "youhou")
	os.Setenv("SYNKER_KAFKA_CACERT", ca.Name())
	os.Setenv("SYNKER_KAFKA_CERT", cert.Name())
	os.Setenv("SYNKER_KAFKA_KEY", key.Name())
	defer os.Unsetenv("SYNKER_KAFKA_SCRAM")
	defer os.Unsetenv("SYNKER_KAFKA_USER")
	defer os.Unsetenv("SYNKER_KAFKA_PASSWORD")
	defer os.Unsetenv("SYNKER_KAFKA_CACERT")
	defer os.Unsetenv("SYNKER_KAFKA_CERT")
	defer os.Unsetenv("SYNKER_KAFKA_KEY")
	_, err = Client()
	assert.Error(err)
}

func TestClient_scram_fail(t *testing.T) {
	assert := assert.New(t)

	os.Setenv("SYNKER_KAFKA_SCRAM", "youhou")
	os.Setenv("SYNKER_KAFKA_USER", "youhou")
	os.Setenv("SYNKER_KAFKA_PASSWORD", "youhou")
	defer os.Unsetenv("SYNKER_KAFKA_SCRAM")
	defer os.Unsetenv("SYNKER_KAFKA_USER")
	defer os.Unsetenv("SYNKER_KAFKA_PASSWORD")
	_, err := Client()
	assert.Error(err)
}

func TestClient_scram_empty(t *testing.T) {
	assert := assert.New(t)

	os.Setenv("SYNKER_KAFKA_USER", "youhou")
	os.Setenv("SYNKER_KAFKA_PASSWORD", "youhou")
	defer os.Unsetenv("SYNKER_KAFKA_USER")
	defer os.Unsetenv("SYNKER_KAFKA_PASSWORD")
	_, err := Client()
	assert.Error(err)
}

func TestCreateTopic(t *testing.T) {
	assert := assert.New(t)

	conn, err := Client()
	assert.Nil(err)
	err = CreateTopic(
		conn,
		models.CreateTopic{
			Name:              test_create_topic,
			NumPartitions:     1,
			ReplicationFactor: 3,
			TopicConfig: []models.TopicConfig{
				{
					Key:   "max.message.bytes",
					Value: "128000",
				},
			},
		})
	assert.Nil(err)
}

func TestListTopics(t *testing.T) {
	assert := assert.New(t)

	conn, err := Client()
	assert.Nil(err)
	_, err = ListTopics(conn)
	assert.Nil(err)
}

func TestProduceMessage(t *testing.T) {
	assert := assert.New(t)

	conn, err := Client()
	assert.Nil(err)
	err = ProduceMessage(
		conn,
		models.KafkaWriteMessage{
			TopicName: test_create_topic,
			Key:       "test",
			Value:     "test",
		})
	assert.Nil(err)
}

func TestConsumeMessage(t *testing.T) {
	assert := assert.New(t)

	conn, err := Client()
	assert.Nil(err)

	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)

	go func() {
		err = consumeMessage(
			conn,
			"test",
			test_create_topic,
		)
		assert.Nil(err)
		<-sigc
		signal.Stop(sigc)
	}()

	err = proc.Signal(os.Interrupt)
	assert.Nil(err)
	time.Sleep(1 * time.Second)
}

func TestDeleteTopics(t *testing.T) {
	assert := assert.New(t)

	conn, err := Client()
	assert.Nil(err)
	err = DeleteTopics(
		conn,
		[]string{
			test_create_topic,
		},
	)
	assert.Nil(err)
}
