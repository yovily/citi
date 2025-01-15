// internal/platform/lookup.go

package platform

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/yovily/customers/citi/auth-service/pkg/resolver"
)

// Add command interface
type commander interface {
	CombinedOutput() ([]byte, error)
}

type commandCreator func(string, ...string) commander

// Update LookupService to include command creator
type LookupService struct {
	resolver    *resolver.Client
	execCommand commandCreator
}

// Update NewLookupService to use real exec.Command by default
func NewLookupService() *LookupService {
	return &LookupService{
		resolver: resolver.NewClient(),
		execCommand: func(name string, args ...string) commander {
			return exec.Command(name, args...)
		},
	}
}

// LookupServer performs platform-specific LDAP server lookup
func (s *LookupService) LookupServer(domain string) (string, error) {
	if domain == "" {
		return "", fmt.Errorf("domain cannot be empty")
	}

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
	case "darwin":
		return s.darwinLookup(domain)
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

func (s *LookupService) linuxLookup(domain string) ([]string, error) {
	cmdStr := fmt.Sprintf("host -t SRV _ldap._tcp.dc._msdcs.%s | awk '{print $NF}'", domain)
	cmd := s.execCommand("sh", "-c", cmdStr)

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
	cmd := s.execCommand("nslookup", "-type=SRV", fmt.Sprintf("_ldap._tcp.dc._msdcs.%s", domain))
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

func (s *LookupService) darwinLookup(domain string) ([]string, error) {
	cmdStr := fmt.Sprintf("host -t SRV _ldap._tcp.dc._msdcs.%s | awk '{print $NF}'", domain)
	cmd := s.execCommand("sh", "-c", cmdStr)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error executing darwin lookup command: %w", err)
	}

	hosts := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(hosts) == 0 {
		return nil, fmt.Errorf("no hosts found")
	}

	return hosts, nil
}
