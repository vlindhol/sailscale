package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vlindhol/sailscale/sailscale-ctl/internal/schema"
)

// TestStatusSchema pins the schema we depend on from `tailscale status --json`.
// If Tailscale changes this output when the bundled tailscale version is bumped,
// this test
// (fed by the fixture, which CI regenerates from the real binary) will fail
// before the change ships.
func TestStatusSchema(t *testing.T) {
	fixture := readFixture(t, "status.json")

	var st schema.DaemonStatus
	if err := json.Unmarshal(fixture, &st); err != nil {
		t.Fatalf("fixture no longer matches DaemonStatus: %v", err)
	}

	if st.BackendState != "Running" {
		t.Errorf("BackendState = %q, want Running", st.BackendState)
	}
	if st.Version != "1.104.7-tabcdef123" {
		t.Errorf("Version = %q", st.Version)
	}
	if !st.TUN {
		t.Error("TUN = false, want true")
	}
	if len(st.TailscaleIPs) != 2 || st.TailscaleIPs[0] != "100.64.0.1" {
		t.Errorf("TailscaleIPs = %v", st.TailscaleIPs)
	}
	if st.Self == nil || st.Self.Name() != "phone.example.ts.net." {
		t.Errorf("Self = %+v", st.Self)
	}
	if got := len(st.Peer); got != 2 {
		t.Errorf("len(Peer) = %d, want 2", got)
	}

	// Rendered view: peers sorted by name, laptop before server.
	view := st.ToView()
	if len(view.Peers) != 2 {
		t.Fatalf("len(view.Peers) = %d, want 2", len(view.Peers))
	}
	if view.Peers[0].Name != "laptop.example.ts.net." {
		t.Errorf("peers[0].Name = %q, want laptop.example.ts.net.", view.Peers[0].Name)
	}
	if view.Peers[1].Name != "server.example.ts.net." {
		t.Errorf("peers[1].Name = %q, want server.example.ts.net.", view.Peers[1].Name)
	}
	if view.Peers[1].Online {
		t.Error("server online = true, want false")
	}
}

// TestVersionSchema pins the schema we depend on from `tailscale version --json`.
func TestVersionSchema(t *testing.T) {
	fixture := readFixture(t, "version.json")

	var v schema.CLIVersion
	if err := json.Unmarshal(fixture, &v); err != nil {
		t.Fatalf("fixture no longer matches CLIVersion: %v", err)
	}
	if v.Long != "1.104.7-tabcdef123" {
		t.Errorf("Long = %q", v.Long)
	}
	if v.Short != "1.104.7" {
		t.Errorf("Short = %q", v.Short)
	}
}

// TestShimExec runs the shim against a fake bundled CLI (a shell script that
// serves the fixtures) and verifies the full exec + parse + emit path.
func TestShimExec(t *testing.T) {
	ctx := context.Background()
	fake := makeFakeCLI(t)

	t.Run("status", func(t *testing.T) {
		var buf bytes.Buffer
		s := newShim(fake)
		if err := s.Status(ctx, &buf); err != nil {
			t.Fatalf("Status: %v", err)
		}
		out := buf.String()
		for _, want := range []string{`"state": "Running"`, `"100.64.0.1"`, `laptop.example.ts.net.`} {
			if !strings.Contains(out, want) {
				t.Errorf("status output missing %q\n%s", want, out)
			}
		}
	})

	t.Run("version", func(t *testing.T) {
		var buf bytes.Buffer
		s := newShim(fake)
		if err := s.Version(ctx, &buf); err != nil {
			t.Fatalf("Version: %v", err)
		}
		out := buf.String()
		if !strings.Contains(out, `1.104.7-tabcdef123`) {
			t.Errorf("version output missing long version\n%s", out)
		}
	})

	t.Run("up and down pass through", func(t *testing.T) {
		s := newShim(fake)
		if err := s.Up(ctx); err != nil {
			t.Errorf("Up: %v", err)
		}
		if err := s.Down(ctx); err != nil {
			t.Errorf("Down: %v", err)
		}
	})
}

// makeFakeCLI writes a fake `tailscale` executable that emits the fixture JSON
// for --json invocations, and records the arguments it was called with.
func makeFakeCLI(t *testing.T) string {
	dir := t.TempDir()
	statusFixture := readFixture(t, "status.json")
	versionFixture := readFixture(t, "version.json")

	script := "#!/bin/sh\n" +
		"case \"$1\" in\n" +
		"  status) cat " + shellQuote(filepath.Join(dir, "status.json")) + " ;;\n" +
		"  version) cat " + shellQuote(filepath.Join(dir, "version.json")) + " ;;\n" +
		"  up|down) exit 0 ;;\n" +
		"esac\n"

	path := filepath.Join(dir, "tailscale")
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "status.json"), statusFixture, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "version.json"), versionFixture, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return b
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}
