// internal/platform/lookup.go

package platform

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/yovily/customers/citi/auth-service/pkg/resolver"
)

// LookupService handles LDAP server lookups for different platforms
type LookupService struct {
	resolver *resolver.Client
}

// NewLookupService creates a new lookup service
func NewLookupService() *LookupService {
	return &LookupService{
		resolver: resolver.NewClient(),
	}
}

// LookupServer performs platform-specific LDAP server lookup
func (s *LookupService) LookupServer(domain string) (string, error) {
	hosts, err := s.getHostsByPlatform(domain)
	if err != nil {
		return "", fmt.Errorf("failed to lookup hosts: %w", err)
	}

	return s.resolver.SelectRandomHost(hosts), nil
}

func (s *LookupService) getHostsByPlatform(domain string) ([]string, error) {
	switch runtime.GOOS {
	case "linux":
		return s.linuxLookup(domain)
	case "windows":
		return s.windowsLookup(domain)
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

func (s *LookupService) linuxLookup(domain string) ([]string, error) {
	cmdStr := fmt.Sprintf("host -t SRV _ldap._tcp.dc._msdcs.%s | awk '{print $NF}'", domain)
	cmd := exec.Command("sh", "-c", cmdStr)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error executing linux lookup command: %w", err)
	}

	hosts := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(hosts) == 0 {
		return nil, fmt.Errorf("no hosts found")
	}

	return hosts, nil
}

func (s *LookupService) windowsLookup(domain string) ([]string, error) {
	cmd := exec.Command("nslookup", "-type=SRV", fmt.Sprintf("_ldap._tcp.dc._msdcs.%s", domain))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error executing windows lookup command: %w", err)
	}

	var hosts []string
	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line, "svr hostname") {
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				hosts = append(hosts, fields[3])
			}
		}
	}

	if len(hosts) == 0 {
		return nil, fmt.Errorf("no hosts found")
	}

	return hosts, nil
}
