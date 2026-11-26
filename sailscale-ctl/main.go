// Command sailscale-ctl is the short-lived Sailscale control binary.
//
// It is a thin shim over the bundled tailscale CLI: it execs the CLI's
// *local* features (up, down, status, version) that talk to tailscaled over
// its LocalAPI, and translates their JSON output for a narrow machine
// contract. It runs on demand and exits; it is not a service.
//
// Boundary: sailscale-ctl owns no Tailscale logic and does not talk to the
// LocalAPI itself — peer negotiation, DERP, WireGuard, routing, and DNS all
// live in tailscaled. Every action here is delegated to the bundled CLI.
//
// Subcommands:
//
//	up       tailscale up
//	down     tailscale down
//	status   tailscale status --json (parsed and re-emitted)
//	version  tailscale version --json (parsed and re-emitted)
package main

import (
	"context"
	"fmt"
	"os"
)

const usage = `sailscale-ctl — drive a local tailscaled via the bundled tailscale CLI

Usage:
  sailscale-ctl [--tailscale PATH] <command>

Commands (local features only):
  up       bring the tailnet online
  down     take the tailnet offline
  status   print self status, Tailscale IPs, and peers (JSON)
  version  print the bundled and client tailscale versions (JSON)

Flags:
  --tailscale PATH   path to the bundled tailscale CLI
                     (default: /usr/bin/sailscale-tailscale)
  -h, --help         show this help
`

func main() {
	ctx := context.Background()

	args := os.Args[1:]
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		fmt.Print(usage)
		os.Exit(0)
	}

	cmd, tailscalePath, err := parseArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sailscale-ctl: %v\n\n%s", err, usage)
		os.Exit(2)
	}

	shim := newShim(tailscalePath)

	switch cmd {
	case "up":
		exitErr(shim.Up(ctx))
	case "down":
		exitErr(shim.Down(ctx))
	case "status":
		exitErr(shim.Status(ctx, os.Stdout))
	case "version":
		exitErr(shim.Version(ctx, os.Stdout))
	case "":
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	default:
		fmt.Fprintf(os.Stderr, "sailscale-ctl: unknown command %q\n\n%s", cmd, usage)
		os.Exit(2)
	}
}

func exitErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "sailscale-ctl: %v\n", err)
		os.Exit(1)
	}
}
