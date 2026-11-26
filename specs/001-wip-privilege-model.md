# Spec 001 — Privilege model (run without root)

## Goal

Let the user install, authenticate, and connect without ever running anything as root, using only the bare-minimum privilege a kernel VPN inherently requires.

## Context

Tailscale's kernel mode needs a `/dev/net/tun` device and `CAP_NET_ADMIN` to create/configure the interface and routes. It does **not** need a root shell, root filesystem access, or UID 0 for its lifetime — it needs only the capability to create a tun interface and add routes.

## Decisions (technical)

- **Use Sailfish's native VPN privilege model** rather than `sudo`/root: the ConnMan plugin's `.conf` runs the helper binary as the unprivileged `nemo` user with the `vpn` group and supplementary capabilities:
  ```ini
  [VPN Binary]
  User = nemo
  Group = vpn
  SupplementaryGroups = inet,net_admin
  ```
- The plugin `.conf` lives at `/etc/connman/vpn-plugin/sailscale.conf` and is installed by the RPM, so there is **no user-facing privilege step**.
- **Rejected as primary:** `--tun=userspace-networking` (zero privilege but TCP/UDP-only, lower throughput, non-transparent proxy routing) and `--operator` mode (needs a persistent privileged bootstrap service).

## Business rules

- Installing Sailscale is the same experience as any other VPN client: no "run as root" step.
- State and secrets live under the unprivileged user's home:
  - `tailscaled` node identity: `/home/nemo/.local/share/sailscale/tailscaled.state`
  - Control socket: `/home/nemo/.local/share/sailscale/tailscaled.sock`
- Interactive web/QR login is preferred; auth keys, if ever used, are single-use and never persisted in plaintext.
- If a target device's kernel lacks `/dev/net/tun`, fall back to userspace-networking as a clearly-labeled degraded mode (pending the kernel `tun` availability check).

## Acceptance criteria

- [ ] `tailscaled` runs as `nemo` (not root) with `vpn` group and `inet,net_admin` only.
- [ ] No root shell, `sudo`, or user-facing elevation is required to install/connect.
- [ ] State/socket files are owned by `nemo` and under the expected paths.
- [ ] Sailjail permissions are declared for all paths and D-Bus services touched.
- [ ] Kernel-tun connectivity works end-to-end with only the `vpn` group granted.

## Out of scope

- Encrypting app preferences (they are plaintext by design; only tailscaled's own state holds keys).
- Supporting devices without `tun` beyond the userspace fallback.