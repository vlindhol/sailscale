# Spec 003 — ConnMan VPN plugin (supervisor + status relay)

## Goal

Make Tailscale appear as a first-class VPN provider in the Sailfish OS networking stack by shipping a ConnMan VPN plugin that supervises `tailscaled` and relays its tunnel state, so the OS connectivity UI shows the VPN and its status exactly like built-in providers.

## Context

Sailfish OS manages all networking through ConnMan. VPNs are provided by shared-library plugins loaded into `connman-vpnd` (the `net.connman.vpn` D-Bus service). A plugin is a thin supervisor between connman-vpnd (network lifecycle) and an external helper binary (the actual tunnel). This matches Tailscale: the plugin runs/stops `tailscaled` and reports its state upward.

## Decisions (technical)

- **Provider identity** is `sailscale`, with provisioning `.config` keys under the `Sailscale.*` namespace.
- **The plugin is a supervisor + status relay**, it does **not** own the tunnel device:
  - On `connect()` it spawns `tailscaled` (with `--state` and `--socket` under `/home/nemo/.local/share/sailscale/`).
  - On `disconnect()` it stops `tailscaled`.
  - It relays Tailscale's self-reported tunnel state into connman-vpnd's `notify()` vocabulary (`IPv4`, routes, `Nameservers`, `ready`/`failure`), so the OS reflects "VPN ready" and the 100.x address.
- **Tailscale owns its own tun and route table.** connman-vpnd does *not* create or hand off a tun interface to Tailscale. The plugin is not involved in tun ownership. (Tailscale is not a push-config client; it computes its own address, routes, and DNS.)
- The plugin reimplements **no** Tailscale logic (no peer negotiation, no DERP, no WireGuard).

## Business rules

- The provider appears as a VPN connection available in the system networking stack.
- Connect/disconnect maps: ConnMan → plugin `connect()`/`disconnect()` → `tailscaled` up/down.
- Tailscale's state is mirrored into `Connection.State` (`idle` → `ready`/`failure`) and emitted via `PropertyChanged`, so the OS shows accurate status and can auto-connect per user preference on network change or roaming.
- Provider properties align with what the controller considers authoritative so the UI and the OS do not disagree about configuration.

## Acceptance criteria

- [ ] `sailscale.so` loads into `connman-vpnd` without error.
- [ ] `connect()` spawns `tailscaled` under `/home/nemo` state/socket paths.
- [ ] `disconnect()` stops `tailscaled`.
- [ ] After a successful connection, connman-vpnd reports `ready` with the 100.x address and correct DNS/route state.
- [ ] System connectivity UI shows the Sailscale VPN entry with correct online/offline status.
- [ ] Plugin does not fight tailscaled over the tun interface.

## Out of scope

- Credential prompting via the ConnMan `Agent` is not used; Tailscale authentication is out-of-band (interactive URL/QR in the UI), separate from the VPN plugin.