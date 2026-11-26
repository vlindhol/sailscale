package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/vlindhol/sailscale/sailscale-ctl/internal/schema"
)

// shim execs the bundled tailscale CLI's local features and re-emits a narrow,
// stable JSON contract. It owns no Tailscale logic; everything is delegated.
type shim struct {
	// tailscalePath is the bundled CLI to exec (overridable for tests).
	tailscalePath string
}

func newShim(path string) *shim {
	return &shim{tailscalePath: path}
}

// Up brings the tailnet online by exec'ing `tailscale up`.
func (s *shim) Up(ctx context.Context) error {
	return s.run(ctx, "up")
}

// Down takes the tailnet offline by exec'ing `tailscale down`.
func (s *shim) Down(ctx context.Context) error {
	return s.run(ctx, "down")
}

// run execs the CLI with the given args, streaming its stdout/stderr to our own
// so the user sees the CLI's messages verbatim. Exit codes propagate as errors.
func (s *shim) run(ctx context.Context, args ...string) error {
	cmd := exec.CommandContext(ctx, s.tailscalePath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tailscale %s: %w", args[0], err)
	}
	return nil
}

// Status execs `tailscale status --json`, parses the daemon's status, and
// re-emits the narrow stable view as JSON to w.
func (s *shim) Status(ctx context.Context, w io.Writer) error {
	out, err := s.output(ctx, "status", "--json")
	if err != nil {
		return err
	}
	var st schema.DaemonStatus
	if err := json.Unmarshal(out, &st); err != nil {
		return fmt.Errorf("parsing `tailscale status --json`: %w", err)
	}
	return writeJSON(w, st.ToView())
}

// Version execs `tailscale version --json` (client version, no daemon required)
// and re-emits a narrow view (bundled tailscale + client versions) as JSON to
// w. The running daemon's version is surfaced by Status instead.
func (s *shim) Version(ctx context.Context, w io.Writer) error {
	out, err := s.output(ctx, "version", "--json")
	if err != nil {
		return err
	}
	var v schema.CLIVersion
	if err := json.Unmarshal(out, &v); err != nil {
		return fmt.Errorf("parsing `tailscale version --json`: %w", err)
	}
	view := schema.VersionView{
		Tailscale: bundledTailscaleVersion,
		Client:    v.Long,
		Daemon:    v.DaemonLong,
	}
	return writeJSON(w, view)
}

// output execs the CLI and captures its stdout for parsing.
func (s *shim) output(ctx context.Context, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, s.tailscalePath, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return nil, fmt.Errorf("tailscale %s: %s", args[0], msg)
		}
		return nil, fmt.Errorf("tailscale %s: %w", args[0], err)
	}
	return out, nil
}

func writeJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// bundledTailscaleVersion is the Tailscale release bundled in the Sailscale
// RPM. It is set at build time via:
//
//	-ldflags "-X main.bundledTailscaleVersion=<version>"
//
// so the bundled version is surfaced independently of the running daemon.
var bundledTailscaleVersion = "unknown"
