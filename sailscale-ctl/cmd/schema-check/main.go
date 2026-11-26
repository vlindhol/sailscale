// Command schema-check validates that JSON produced by the real bundled
// tailscale CLI still decodes with the same types sailscale-ctl uses.
//
// It is used by CI after a bundled tailscale bump: run the actual `tailscale
// status --json` and `tailscale version --json`, then feed both files here.
// A non-zero exit means the CLI's schema changed and the translation layer
// (internal/schema) needs updating.
//
// Usage: schema-check STATUS_JSON VERSION_JSON
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/vlindhol/sailscale/sailscale-ctl/internal/schema"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: schema-check STATUS_JSON VERSION_JSON")
		os.Exit(2)
	}

	if err := checkStatus(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "status schema check failed: %v\n", err)
		os.Exit(1)
	}
	if err := checkVersion(os.Args[2]); err != nil {
		fmt.Fprintf(os.Stderr, "version schema check failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("schema-check: OK")
}

func checkStatus(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var st schema.DaemonStatus
	if err := json.Unmarshal(b, &st); err != nil {
		return fmt.Errorf("`status --json` no longer decodes: %w", err)
	}
	// Key fields must be present and populated in a live, authenticated or
	// NeedsLogin daemon.
	if st.BackendState == "" {
		return fmt.Errorf("BackendState is empty")
	}
	if st.Version == "" {
		return fmt.Errorf("Version is empty")
	}
	_ = st.ToView() // ensure the projection still runs without panic
	return nil
}

func checkVersion(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var v schema.CLIVersion
	if err := json.Unmarshal(b, &v); err != nil {
		return fmt.Errorf("`version --json` no longer decodes: %w", err)
	}
	if v.Long == "" {
		return fmt.Errorf("long version is empty")
	}
	return nil
}
