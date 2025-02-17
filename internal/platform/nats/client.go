package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

type Config struct {
	URL           string
	MaxReconnects int
	ReconnectWait time.Duration
}

type Client struct {
	conn *nats.Conn
}

func NewClient(config Config) (*Client, error) {
	opts := []nats.Option{
		nats.MaxReconnects(config.MaxReconnects),
		nats.ReconnectWait(config.ReconnectWait),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			fmt.Printf("Disconnected from NATS: %v\n", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			fmt.Printf("Reconnected to NATS server %v\n", nc.ConnectedUrl())
		}),
	}

	nc, err := nats.Connect(config.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	return &Client{conn: nc}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Publish(subject string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := c.conn.Publish(subject, payload); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (c *Client) Subscribe(subject string, handler func([]byte) error) (*nats.Subscription, error) {
	sub, err := c.conn.Subscribe(subject, func(msg *nats.Msg) {
		if err := handler(msg.Data); err != nil {
			fmt.Printf("Error handling message: %v\n", err)
		}
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	return sub, nil
}

func (c *Client) Request(subject string, data interface{}, timeout time.Duration) ([]byte, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	msg, err := c.conn.Request(subject, payload, timeout)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return msg.Data, nil
}

func (c *Client) QueueSubscribe(subject, queue string, handler func([]byte) error) (*nats.Subscription, error) {
	sub, err := c.conn.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		if err := handler(msg.Data); err != nil {
			fmt.Printf("Error handling message: %v\n", err)
		}
	})
	if err != nil {
		return nil, fmt.Errorf("failed to queue subscribe: %w", err)
	}

	return sub, nil
}
