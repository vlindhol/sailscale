// Package sailscale holds shared constants for the sailscale-ctl binary.
package sailscale

// TailscaleCLI is the absolute path where the RPM installs the bundled
// tailscale CLI (renamed to sailscale-tailscale to avoid clashing with any
// system binary). sailscale-ctl execs this for tailscale's local features.
const TailscaleCLI = "/usr/bin/sailscale-tailscale"
