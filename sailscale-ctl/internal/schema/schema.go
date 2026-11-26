// Package schema pins the subset of the bundled tailscale CLI's JSON output
// that sailscale-ctl depends on, and the narrow, stable JSON it re-emits.
//
// The CLI's --json output is documented as "subject to change"; this package
// is the single place that knows those shapes, so a CI job can decode real
// CLI output with the same code sailscale-ctl uses and fail the build if a
// schema change would break it.
package schema

import "sort"

// StatusView is the stable JSON sailscale-ctl emits for a status command.
type StatusView struct {
	State        string     `json:"state"`
	Version      string     `json:"version"`
	TUN          bool       `json:"tun"`
	TailscaleIPs []string   `json:"tailscaleIPs,omitempty"`
	Self         *PeerView  `json:"self,omitempty"`
	Peers        []PeerView `json:"peers"`
}

// PeerView is a single peer as sailscale-ctl surfaces it.
type PeerView struct {
	Name   string   `json:"name"`
	IPs    []string `json:"ips,omitempty"`
	Online bool     `json:"online"`
}

// VersionView is the stable JSON sailscale-ctl emits for a version command.
type VersionView struct {
	Tailscale string `json:"tailscale"`
	Client    string `json:"client,omitempty"`
	Daemon    string `json:"daemon,omitempty"`
}

// DaemonStatus is the subset of the CLI's `status --json` (ipnstate.Status)
// that sailscale-ctl reads. Field names match the CLI's JSON exactly.
type DaemonStatus struct {
	Version      string                `json:"Version"`
	BackendState string                `json:"BackendState"`
	TUN          bool                  `json:"TUN"`
	TailscaleIPs []string              `json:"TailscaleIPs"`
	Self         *PeerStatus           `json:"Self"`
	Peer         map[string]PeerStatus `json:"Peer"`
}

// ToView projects the daemon status into the narrow, stable status view.
func (d DaemonStatus) ToView() StatusView {
	v := StatusView{
		State:        d.BackendState,
		Version:      d.Version,
		TUN:          d.TUN,
		TailscaleIPs: d.TailscaleIPs,
	}
	if d.Self != nil {
		v.Self = &PeerView{Name: d.Self.Name(), IPs: d.Self.TailscaleIPs, Online: d.Self.Online}
	}
	v.Peers = peerViews(d.Peer)
	return v
}

// PeerStatus is the subset of a peer's status sailscale-ctl reads.
type PeerStatus struct {
	DNSName      string   `json:"DNSName"`
	HostName     string   `json:"HostName"`
	TailscaleIPs []string `json:"TailscaleIPs"`
	Online       bool     `json:"Online"`
}

// Name returns the peer's DNS name, falling back to its host name.
func (p PeerStatus) Name() string {
	if p.DNSName != "" {
		return p.DNSName
	}
	return p.HostName
}

func peerViews(m map[string]PeerStatus) []PeerView {
	out := make([]PeerView, 0, len(m))
	for _, p := range m {
		out = append(out, PeerView{Name: p.Name(), IPs: p.TailscaleIPs, Online: p.Online})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// CLIVersion is the subset of `tailscale version --json` (version.Meta) that
// sailscale-ctl reads.
type CLIVersion struct {
	Long       string `json:"long"`
	Short      string `json:"short"`
	DaemonLong string `json:"daemonLong,omitempty"`
}
