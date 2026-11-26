package main

import (
	"testing"

	c "github.com/vlindhol/sailscale/sailscale-ctl/internal/sailscale"
)

func TestParseArgs(t *testing.T) {
	t.Run("no args", func(t *testing.T) {
		cmd, cli, err := parseArgs(nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cmd != "" {
			t.Errorf("cmd = %q, want empty", cmd)
		}
		if cli != c.TailscaleCLI {
			t.Errorf("tailscale = %q, want default %q", cli, c.TailscaleCLI)
		}
	})

	t.Run("command only", func(t *testing.T) {
		cmd, _, err := parseArgs([]string{"status"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cmd != "status" {
			t.Errorf("cmd = %q, want %q", cmd, "status")
		}
	})

	t.Run("tailscale flag before command", func(t *testing.T) {
		cmd, cli, err := parseArgs([]string{"--tailscale", "/tmp/tailscale", "status"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cmd != "status" {
			t.Errorf("cmd = %q, want %q", cmd, "status")
		}
		if cli != "/tmp/tailscale" {
			t.Errorf("tailscale = %q, want /tmp/tailscale", cli)
		}
	})

	t.Run("tailscale flag after command", func(t *testing.T) {
		cmd, cli, err := parseArgs([]string{"up", "--tailscale", "/tmp/tailscale"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cmd != "up" {
			t.Errorf("cmd = %q, want %q", cmd, "up")
		}
		if cli != "/tmp/tailscale" {
			t.Errorf("tailscale = %q, want /tmp/tailscale", cli)
		}
	})

	t.Run("tailscale equals form", func(t *testing.T) {
		_, cli, err := parseArgs([]string{"--tailscale=/tmp/tailscale", "down"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cli != "/tmp/tailscale" {
			t.Errorf("tailscale = %q, want /tmp/tailscale", cli)
		}
	})

	t.Run("--tailscale with no value", func(t *testing.T) {
		_, _, err := parseArgs([]string{"--tailscale"})
		if err == nil {
			t.Fatal("expected error for --tailscale without a path")
		}
	})

	t.Run("extra positional arg", func(t *testing.T) {
		_, _, err := parseArgs([]string{"status", "extra"})
		if err == nil {
			t.Fatal("expected error for unexpected positional argument")
		}
	})
}
