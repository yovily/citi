// pkg/resolver/resolver_test.go

package resolver

import (
    "testing"
)

func TestSelectRandomHost(t *testing.T) {
    tests := []struct {
        name     string
        hosts    []string
        wantHost bool // true if we expect a host, false if we expect empty string
    }{
        {
            name:     "empty hosts list",
            hosts:    []string{},
            wantHost: false,
        },
        {
            name:     "single host",
            hosts:    []string{"host1"},
            wantHost: true,
        },
        {
            name:     "multiple hosts",
            hosts:    []string{"host1", "host2", "host3"},
            wantHost: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            client := NewClient()
            got := client.SelectRandomHost(tt.hosts)

            if tt.wantHost {
                // Check if returned host is in the input slice
                found := false
                for _, h := range tt.hosts {
                    if got == h {
                        found = true
                        break
                    }
                }
                if !found {
                    t.Errorf("SelectRandomHost() = %v, want one of %v", got, tt.hosts)
                }
            } else {
                if got != "" {
                    t.Errorf("SelectRandomHost() = %v, want empty string", got)
                }
            }
        })
    }
}

func TestRandomness(t *testing.T) {
    hosts := []string{"host1", "host2", "host3", "host4", "host5"}
    client := NewClient()
    
    // Track frequency of each host
    frequency := make(map[string]int)
    iterations := 1000

    for i := 0; i < iterations; i++ {
        host := client.SelectRandomHost(hosts)
        frequency[host]++
    }

    // Check that each host was selected at least once
    for _, h := range hosts {
        if frequency[h] == 0 {
            t.Errorf("Host %s was never selected in %d iterations", h, iterations)
        }
    }
}