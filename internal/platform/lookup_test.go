// internal/platform/lookup_test.go

package platform

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/yovily/citi/internal/platform/resolver"
)

// Mock command executor
type mockCmd struct {
	output []byte
	err    error
}

func (m *mockCmd) CombinedOutput() ([]byte, error) {
	return m.output, m.err
}

// Mock command creator
func mockExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestLookupService(t *testing.T) {
	service := NewLookupService()
	if service == nil {
		t.Fatal("Expected non-nil service")
	}
}

func TestLookupServer(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		mockOutput []byte
		mockError  error
		wantError  bool
	}{
		{
			name:       "valid domain",
			domain:     "example.com",
			mockOutput: []byte("ldap1.example.com.\nldap2.example.com."),
			mockError:  nil,
			wantError:  false,
		},
		{
			name:      "empty domain",
			domain:    "",
			mockError: fmt.Errorf("domain cannot be empty"),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock command
			mockCmd := &mockCmd{
				output: tt.mockOutput,
				err:    tt.mockError,
			}

			// Create service with mock command
			svc := &LookupService{
				resolver: resolver.NewClient(),
				execCommand: func(name string, args ...string) commander {
					return mockCmd
				},
			}

			// Run the test
			host, err := svc.LookupServer(tt.domain)

			if (err != nil) != tt.wantError {
				t.Errorf("LookupServer() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError && host == "" {
				t.Error("LookupServer() returned empty host when error not expected")
			}
		})
	}
}

// Helper process to mock command execution
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(0)
}
