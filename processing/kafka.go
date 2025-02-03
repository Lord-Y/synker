// Package processing provide all requirements to process change data capture
package processing

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
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
)

const (
	timeout time.Duration = 10 * time.Second
)

// kClient permit to connect to kafka brokers
func (c *Validate) kClient() (conn *kafka.Conn, err error) {
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

func (c *Validate) connectToController(conn *kafka.Conn) (connLeader *kafka.Conn, err error) {
	controller, err := conn.Controller()
	if err != nil {
		return
	}

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
	return
}

// createTopic permit to create a topic
func (c *Validate) createTopic(conn *kafka.Conn, kf createTopicModel) (err error) {
	defer conn.Close()
	connLeader, err := c.connectToController(conn)
	if err != nil {
		return
	}
	defer connLeader.Close()

	var topicConfig []kafka.ConfigEntry
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

// listTopics permit to list all topics
func (c *Validate) listTopics(conn *kafka.Conn) (topics []string, err error) {
	defer conn.Close()
	connLeader, err := c.connectToController(conn)
	if err != nil {
		return
	}
	defer connLeader.Close()

	partitions, err := connLeader.ReadPartitions()
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

// deleteTopics permit to delete topics
func (c *Validate) deleteTopics(conn *kafka.Conn, topics []string) (err error) {
	defer conn.Close()
	connLeader, err := c.connectToController(conn)
	if err != nil {
		c.Logger.Info().Msgf("XXXX %v", err.Error())
		return
	}
	defer connLeader.Close()

	err = connLeader.DeleteTopics(topics...)
	return
}

// produceMessage permit to write a message into specified topic
func (c *Validate) produceMessage(conn *kafka.Conn, message kafkaWriteMessage) (err error) {
	defer conn.Close()
	connLeader, err := c.connectToController(conn)
	if err != nil {
		return
	}
	defer connLeader.Close()

	err = connLeader.SetWriteDeadline(time.Now().Add(10 * time.Second))
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
func (c *Validate) consumeMessage(conn *kafka.Conn, consumerGroup string, topicName string) (err error) {
	defer conn.Close()
	connLeader, err := c.connectToController(conn)
	if err != nil {
		c.Logger.Info().Msgf("YYYY %v", err.Error())
		return
	}
	defer connLeader.Close()

	brokers := []string{
		connLeader.Broker().Host + ":" + strconv.Itoa(connLeader.Broker().Port),
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
	c.Logger.Debug().Msgf("Message at offset %d: %s = %s", m.Offset, string(m.Key), string(m.Value))
	return
}
