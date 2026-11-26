# sailscale-ctl

The short-lived Go control binary that drives the bundled `tailscaled`.

## Role

- Runs on demand and exits; it is **not** a service.
- It is a **thin shim over the bundled `tailscale` CLI**: every command execs
  the CLI's *local* features and translates their output. It does not talk to
  the LocalAPI itself and owns no Tailscale logic.

## Command → CLI mapping

| Command | Delegates to |
|---------|--------------|
| `up` | `tailscale up` |
| `down` | `tailscale down` |
| `status` | `tailscale status --json` (parsed, re-emitted) |
| `version` | `tailscale version --json` (parsed, re-emitted) |

## Business rules

- The bundled CLI is located at `--tailscale PATH` (default
  `/usr/bin/sailscale-tailscale`, the RPM-installed rename).
- `status`/`version` emit stable JSON; only the fields Sailscale surfaces are
  kept, so consumers do not depend on the CLI's "subject to change" text
  output.
- The CLI's JSON schemas live under `testdata/` and are pinned by tests; CI
  re-verifies them whenever the bundled Tailscale version changes.

## Boundary rule

`sailscale-ctl` owns **no** Tailscale logic — no peer negotiation, DERP,
WireGuard, routing, or DNS. All of that lives in the bundled `tailscaled`.
Every action here is delegated to the bundled CLI.