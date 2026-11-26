package sailscale

import "testing"

func TestTailscaleCLI(t *testing.T) {
	if TailscaleCLI != "/usr/bin/sailscale-tailscale" {
		t.Errorf("TailscaleCLI = %q, want /usr/bin/sailscale-tailscale", TailscaleCLI)
	}
}
