package surang

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Surang represents an SSH tunnel.
type Surang struct {
	Name     string
	Command  string
	ExpectIP string
	Port     int
}

// IsRunning checks if the surang is currently running.
func (s *Surang) IsRunning(ctx context.Context) bool {
	cmdPattern := fmt.Sprintf("ssh -f -N -D %d %s", s.Port, s.Command)
	cmd := exec.CommandContext(ctx, "pgrep", "-f", cmdPattern)
	output, err := cmd.Output()
	return err == nil && len(output) > 0
}

// Start launches the surang.
func (s *Surang) Start(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "ssh", "-f", "-N", "-D", strconv.Itoa(s.Port), s.Command)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting surang %s: %w", s.Name, err)
	}
	return nil
}

// Restart restarts the surang.
func (s *Surang) Restart(ctx context.Context) error {
	if err := s.Stop(ctx); err != nil {
		return err
	}
	return s.Start(ctx)
}

// Stop terminates the surang.
func (s *Surang) Stop(ctx context.Context) error {
	cmdPattern := fmt.Sprintf("ssh -f -N -D %d %s", s.Port, s.Command)
	cmd := exec.CommandContext(ctx, "pkill", "-f", cmdPattern)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error stopping surang %s: %w", s.Name, err)
	}
	return nil
}

// Check verifies if the surang is correctly forwarding traffic.
func (s *Surang) Check(ctx context.Context) (bool, error) {
	cmd := exec.CommandContext(ctx, "curl", "-4", "-s", "--socks5", fmt.Sprintf("localhost:%d", s.Port), "icanhazip.com")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("error checking IP for surang %s: %w", s.Name, err)
	}
	return strings.TrimSpace(string(output)) == s.ExpectIP, nil
}
