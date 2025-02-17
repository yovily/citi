package resolver

import (
	"math/rand"
	"time"
)

type Client struct {
	rand *rand.Rand
}

func NewClient() *Client {
	return &Client{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (c *Client) SelectRandomHost(hosts []string) string {
	if len(hosts) == 0 {
		return ""
	}
	return hosts[c.rand.Intn(len(hosts))]
}
