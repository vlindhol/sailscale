package main

import (
	"fmt"

	c "github.com/vlindhol/sailscale/sailscale-ctl/internal/sailscale"
)

// parseArgs splits the argument list into the command and the bundled tailscale
// CLI path (from --tailscale/--tailscale=PATH). It returns an error when a flag
// is malformed, and "" for cmd when no command is present.
func parseArgs(args []string) (cmd, tailscalePath string, err error) {
	tailscalePath = c.TailscaleCLI
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--tailscale":
			if i+1 >= len(args) {
				return "", "", fmt.Errorf("--tailscale requires a path")
			}
			tailscalePath = args[i+1]
			i++
		case len(a) >= len("--tailscale=") && a[:len("--tailscale=")] == "--tailscale=":
			tailscalePath = a[len("--tailscale="):]
		default:
			if cmd == "" {
				cmd = a
			} else {
				return "", "", fmt.Errorf("unexpected argument %q", a)
			}
		}
	}
	return cmd, tailscalePath, nil
}
