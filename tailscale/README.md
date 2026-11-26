# Tailscale bundle

The Sailscale RPM bundles the official Tailscale rather than reimplementing it.
`tailscale/` records the pinned build that gets vendored into the RPM; it is
**not** our source.

The pinned version is stored in `tailscale/VERSION` — the **single source of
truth**. The SHA-256 below and `tailscale/fetch.sh` must stay in agreement with
it; CI verifies this.

## Pin

| Field | Value |
|-------|-------|
| Track | `stable` |
| Version | `1.104.7` |
| Architecture | `arm64` (aarch64 — Jolla C2 and later; Sailfish OS 5.x) |
| Tarball | `tailscale_1.104.7_arm64.tgz` |
| URL | `https://pkgs.tailscale.com/stable/tailscale_1.104.7_arm64.tgz` |
| SHA-256 | `d6bed2720310b328b861f943d7b4016c5358a9bad60575eb51e8143b79068484` |

The tarball contains two statically linked binaries — `tailscaled` (the resident
daemon that Sailscale supervises) and `tailscale` (the CLI, which `sailscale-ctl`
execs for its local features).

## Updating the pin

1. Choose the newest version published on the **stable** track
   (`https://pkgs.tailscale.com/stable/`).
2. Download the `arm64` tarball and its `.sha256`.
3. Update `tailscale/VERSION`, then this table (version, tarball name, SHA-256)
   and `tailscale/fetch.sh`'s checksum.
4. Bump the places that reference the version so CI's consistency check stays
   green:
   - `sailscale-ctl` build's `-X main.bundledTailscaleVersion=<version>`.
   - `packaging/sailscale.spec`'s `%define tailscale_version`.

CI's `verify-tailscale-consistency` job fails if these disagree with
`tailscale/VERSION`.