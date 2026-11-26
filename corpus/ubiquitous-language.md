# Ubiquitous language

The agreed-upon vocabulary for Sailscale. Every new term below is the
authoritative name for a thing that exists in the code; when you write code,
corpus, or change descriptions, use these exact words and nothing else. Do not
coin near-synonyms (for example, do not call Tailscale "the engine").

## Product & units

| Term | What it names |
|------|---------------|
| **Sailscale** | This project: a Tailscale client for Sailfish OS. |
| **Tailscale** | The upstream product Sailscale wraps. Never reimplemented; always bundled. |
| **`sailscale-ctl`** | The short-lived Go control binary. A backend-for-frontend over the CLI. |
| **`tailscaled`** | The resident Tailscale daemon that holds node state and drives the tunnel. |
| **`tailscale` (CLI)** | The bundled Tailscale command-line client, renamed `sailscale-tailscale` on install. |
| **backend-for-frontend** (BFF) | What `sailscale-ctl` is: a narrow, stable interface in front of the CLI. |

## Directories

| Term | What it names |
|------|---------------|
| **`tailscale/`** | Holds the pinned Tailscale version (`VERSION`) and the fetch script. |
| **`sailscale-ctl/`** | The Go module for the control binary. |
| **`packaging/`** | The RPM spec and the ConnMan plugin `.conf`. |
| **`corpus/`** | Descriptions of the current state of the code (where this file lives). |
| **`specs/`** | Forward-looking work documents (not referenced by the corpus). |

## Runtime concepts

| Term | What it names |
|------|---------------|
| **LocalAPI** | Tailscale's local control interface exposed by `tailscaled`; the CLI speaks it. `sailscale-ctl` does **not** — it delegates to the CLI. |
| **tailnet** | A Tailscale network that a device joins. |
| **command set** | The closed set of `sailscale-ctl` commands: `up`, `down`, `status`, `version`. |
| **translation layer** | `internal/schema` — the single place that knows the CLI's `--json` shapes. |

## Rules

- Code uses these names (directories, binaries, types, fields) as-is. Do not
  re-describe a thing with a new metaphor.
- A term that names no code yet (for example, the ConnMan plugin) is **not**
  part of the ubiquitous language until its code exists.