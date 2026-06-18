package service

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
)

type TransactionProducer interface {
	PublishTransaction(ctx context.Context, transactionId int) error
}

type NATSTransactionProducer struct {
	conn *nats.Conn
}

func NewNATSTransactionProducer(conn *nats.Conn) TransactionProducer {
	return &NATSTransactionProducer{conn: conn}
}

type transactionMessage struct {
	ID    int    `json:"id"`
	Event string `json:"event"`
}

func (p *NATSTransactionProducer) PublishTransaction(ctx context.Context, transactionId int) error {
	msg := transactionMessage{
		ID:    transactionId,
		Event: "created",
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.conn.Publish("transactions.created", data)
}
