package kafka

import (
	"crypto/tls"
	kafgo "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"serviceB/cfg"
	"serviceB/internal/services"
	"time"
)

func NewConsumer(cfg *cfg.KafkaConfig, service services.IUserBalance) (Consumer, error) {
	mechanism, err := scram.Mechanism(scram.SHA256, cfg.SASL.User, cfg.SASL.Password)
	if err != nil {
		return Consumer{}, err
	}
	dialer := &kafgo.Dialer{
		Timeout:       cfg.DialerTimeout * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
		TLS: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	reader := kafgo.NewReader(kafgo.ReaderConfig{
		Brokers:  cfg.BrokerURLs,
		GroupID:  cfg.ConsumerGroup,
		Topic:    cfg.Topic,
		MinBytes: cfg.MinBytes,
		MaxBytes: cfg.MaxBytes,
		Dialer:   dialer,
	})

	return Consumer{
		cfg:                cfg,
		reader:             reader,
		userBalanceService: service,
	}, nil
}
