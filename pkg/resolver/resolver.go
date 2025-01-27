// pkg/resolver/resolver.go

package resolver

import (
    "math/rand"
    "time"
)

// Client handles LDAP server resolution
type Client struct {
    r *rand.Rand
}

// NewClient creates a new resolver client
func NewClient() *Client {
    return &Client{
        // In Go 1.21+, we use a local random source instead of global rand.Seed
        r: rand.New(rand.NewSource(time.Now().UnixNano())),
    }
}

// SelectRandomHost picks a random host from the available servers
func (c *Client) SelectRandomHost(hosts []string) string {
    if len(hosts) == 0 {
        return ""
    }
    return hosts[c.r.Intn(len(hosts))]
}