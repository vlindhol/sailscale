package schema

import (
	"encoding/json"
	"testing"
)

// These tests pin the single translation layer between the tailscale CLI's
// JSON and sailscale-ctl's stable output. If the CLI changes shape on a version
// bump, these are what break — and they break *here only*, not in the BFF's
// consumers.

func TestDaemonStatusDecodeAndView(t *testing.T) {
	// A minimal subset shaped like real `tailscale status --json`.
	in := []byte(`{
		"Version": "1.104.7-tabc",
		"TUN": true,
		"BackendState": "Running",
		"TailscaleIPs": ["100.64.0.1"],
		"Self": {
			"HostName": "phone",
			"DNSName": "phone.example.ts.net.",
			"TailscaleIPs": ["100.64.0.1"],
			"Online": true
		},
		"Peer": {
			"key2": {
				"HostName": "server",
				"DNSName": "server.example.ts.net.",
				"TailscaleIPs": ["100.64.0.3"],
				"Online": false
			},
			"key1": {
				"HostName": "laptop",
				"DNSName": "laptop.example.ts.net.",
				"TailscaleIPs": ["100.64.0.2"],
				"Online": true
			}
		}
	}`)

	var st DaemonStatus
	if err := json.Unmarshal(in, &st); err != nil {
		t.Fatalf("decode: %v", err)
	}

	view := st.ToView()
	if view.State != "Running" {
		t.Errorf("view.State = %q", view.State)
	}
	if view.Self == nil || view.Self.Name != "phone.example.ts.net." {
		t.Errorf("view.Self = %+v", view.Self)
	}
	if len(view.Peers) != 2 {
		t.Fatalf("len(view.Peers) = %d, want 2", len(view.Peers))
	}
	// Peers sorted by name: laptop < server.
	if view.Peers[0].Name != "laptop.example.ts.net." {
		t.Errorf("peers[0].Name = %q", view.Peers[0].Name)
	}
	if view.Peers[1].Name != "server.example.ts.net." {
		t.Errorf("peers[1].Name = %q", view.Peers[1].Name)
	}
	if view.Peers[1].Online {
		t.Error("server should be offline")
	}
}

func TestPeerNameFallback(t *testing.T) {
	if got := (PeerStatus{HostName: "host"}).Name(); got != "host" {
		t.Errorf("Name () = %q, want host", got)
	}
	if got := (PeerStatus{DNSName: "host.example.ts.net.", HostName: "host"}).Name(); got != "host.example.ts.net." {
		t.Errorf("Name() = %q, want DNS name", got)
	}
}

func TestCLIVersionDecode(t *testing.T) {
	in := []byte(`{"long": "1.104.7-tabc", "short": "1.104.7", "daemonLong": "1.104.7-tabc"}`)
	var v CLIVersion
	if err := json.Unmarshal(in, &v); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if v.Long != "1.104.7-tabc" || v.Short != "1.104.7" {
		t.Errorf("decode mismatch: %+v", v)
	}
}
