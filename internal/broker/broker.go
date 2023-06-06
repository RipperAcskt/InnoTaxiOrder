package broker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/RipperAcskt/innotaxiorder/config"
	"github.com/RipperAcskt/innotaxiorder/internal/model"
	"github.com/segmentio/kafka-go"
)

type Broker struct {
	conn *kafka.Conn
	cfg  *config.Config
}

func New(cfg *config.Config) (*Broker, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", cfg.BROKER_HOST, "order", 0)
	if err != nil {
		return nil, fmt.Errorf("dial leader failed: %w", err)
	}

	return &Broker{
		conn: conn,
		cfg:  cfg,
	}, nil
}

func (b *Broker) Write(order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("marshal failed: %w", err)
	}

	_, err = b.conn.Write(data)
	if err != nil {
		return fmt.Errorf("write failed: %w", err)
	}
	return nil
}
