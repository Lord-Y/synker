// Kafka package will handle all kafka requirements
package kafka

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Lord-Y/synker/commons"
	"github.com/Lord-Y/synker/models"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
)

const (
	timeout time.Duration = 10 * time.Second
)

// Client permit to connect to kafka brokers
func Client() (conn *kafka.Conn, err error) {
	var (
		dialer    *kafka.Dialer
		mechanism sasl.Mechanism
	)

	switch {
	case commons.GetKafkaCACert() != "" &&
		commons.GetKafkaCert() != "" &&
		commons.GetKafkaKey() != "" &&
		commons.GetKafkaUser() != "" &&
		commons.GetKafkaPassword() != "":
		var (
			cert      tls.Certificate
			tlsConfig tls.Config
			ca        []byte
		)
		cert, err = tls.LoadX509KeyPair(commons.GetKafkaCert(), commons.GetKafkaKey())
		if err != nil {
			return
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
		ca, err = os.ReadFile(commons.GetKafkaCACert())
		if err != nil {
			return
		}
		caPool := x509.NewCertPool()
		caPool.AppendCertsFromPEM(ca)

		tlsConfig.RootCAs = caPool

		dialer = &kafka.Dialer{
			Timeout:   timeout,
			DualStack: true,
			TLS:       &tlsConfig,
		}
	case commons.GetKafkaScram() != "" &&
		commons.GetKafkaUser() != "" &&
		commons.GetKafkaPassword() != "":
		mechanism, err = scram.Mechanism(
			scram.SHA256,
			commons.GetKafkaUser(),
			commons.GetKafkaPassword(),
		)
		if err != nil {
			return
		}
		dialer = &kafka.Dialer{
			Timeout:       timeout,
			DualStack:     true,
			SASLMechanism: mechanism,
		}
	case commons.GetKafkaScram() == "" &&
		commons.GetKafkaUser() != "" &&
		commons.GetKafkaPassword() != "":
		mechanism := plain.Mechanism{
			Username: commons.GetKafkaUser(),
			Password: commons.GetKafkaPassword(),
		}
		dialer = &kafka.Dialer{
			Timeout:       timeout,
			DualStack:     true,
			SASLMechanism: mechanism,
		}
	default:
		dialer = &kafka.Dialer{
			Timeout:   timeout,
			DualStack: true,
		}
	}
	conn, err = dialer.Dial("tcp", commons.GetKafkaURI())

	if err != nil {
		return
	}
	return
}

// CreateTopic permit to create a topic
func CreateTopic(conn *kafka.Conn, kf models.CreateTopic) (err error) {
	defer conn.Close()
	controller, err := conn.Controller()
	if err != nil {
		return
	}

	var (
		connLeader  *kafka.Conn
		topicConfig []kafka.ConfigEntry
	)
	connLeader, err = kafka.Dial(
		"tcp",
		net.JoinHostPort(
			controller.Host,
			strconv.Itoa(controller.Port),
		),
	)
	if err != nil {
		return
	}
	defer connLeader.Close()
	if len(kf.TopicConfig) > 0 {
		for _, v := range kf.TopicConfig {
			topicConfig = append(
				topicConfig,
				kafka.ConfigEntry{
					ConfigName:  v.Key,
					ConfigValue: v.Value,
				},
			)
		}
	}

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             kf.Name,
			NumPartitions:     kf.NumPartitions,
			ReplicationFactor: kf.ReplicationFactor,
			ConfigEntries:     topicConfig,
		},
	}

	err = connLeader.CreateTopics(topicConfigs...)
	if err != nil {
		return
	}
	return
}

// ListTopics permit to list all topics
func ListTopics(conn *kafka.Conn) (topics []string, err error) {
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return
	}

	m := map[string]struct{}{}
	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}
	for k := range m {
		topics = append(topics, k)
	}
	return
}

// DeleteTopics permit to delete topics
func DeleteTopics(conn *kafka.Conn, topics []string) (err error) {
	defer conn.Close()

	err = conn.DeleteTopics(topics...)
	return
}

// ProduceMessage permit to write a message into specified topic
func ProduceMessage(conn *kafka.Conn, message models.KafkaWriteMessage) (err error) {
	defer conn.Close()

	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return
	}
	_, err = conn.WriteMessages(
		kafka.Message{
			Key:   []byte(message.Key),
			Value: []byte(message.Value),
		},
	)

	w := &kafka.Writer{
		Addr: kafka.TCP(
			strings.Join(
				[]string{
					conn.Broker().Host,
					strconv.Itoa(conn.Broker().Port),
				},
				":",
			),
		),
		Topic:    message.TopicName,
		Balancer: &kafka.LeastBytes{},
	}

	err = w.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte(message.Key),
			Value: []byte(message.Value),
		},
	)
	return
}

// consumeMessage permit to consume message into specified topic
// and will be used for unit testing only
func consumeMessage(conn *kafka.Conn, consumerGroup string, topicName string) (err error) {
	defer conn.Close()
	brokers := []string{
		conn.Broker().Host,
		strconv.Itoa(conn.Broker().Port),
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topicName,
		GroupID:  consumerGroup,
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	defer r.Close()

	ctx := context.Background()
	m, err := r.ReadMessage(ctx)
	if err != nil {
		return
	}
	log.Debug().Msgf("Message at offset %d: %s = %s", m.Offset, string(m.Key), string(m.Value))
	return
}
