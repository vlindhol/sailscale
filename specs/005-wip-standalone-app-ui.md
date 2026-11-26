# Spec 005 — Standalone app UI (status, peers, settings)

## Goal

Provide a Sailfish-native QML/Silica app with a status/home screen, a login screen, and settings, giving parity with what the `tailscale` CLI and the device-relevant parts of the admin console can do, within reason.

## Context

Sailscale has two UI surfaces: the standalone app, and an entry in the general Settings app. The standalone app is the primary surface for advanced configuration.

## Decisions (technical)

- Built with Silica components (`Page`, `Label`, `Button`, `ComboBox`, `TextField`, `Switch`, `BusyIndicator`, pulley menus).
- Follows Sailfish navigation/pulley/gesture conventions.
- App branding: name **Sailscale**, a **sail motif** icon. The word "Tailscale" may appear in copy, but Tailscale's logo is not reused and no first-party status is claimed.

## Business rules

- **Status/home** shows: online state, IP (100.x), key-expiry countdown, self-status, and a peer list (names, IPs, online/reachability).
- **Settings** exposes: hostname, exit-node picker, accept-routes toggle + route list, MagicDNS toggle (default **on**), key re-auth/reset, and logout (with optional wipe).
- **Advanced** exposes: networking-mode display, engine version/update info, and placeholders for deferred server-role/SSH features.
- Exit-node behavior: default **no exit node**, but remember the last-chosen one and offer a quick "use last exit node" action.
- Primary CTA reflects state: **Log in** → **Connecting** → **Connected** (with IP) → **Disconnect**.
- Errors (auth timeout, network change, key expiry) surface **inline** and, for key events, via a Sailfish notification.
- Roaming/network change re-establishes the tunnel silently; status reflects "reconnecting" without user action.

## Acceptance criteria

- [ ] Home screen shows online state, IP, expiry, and the peer list.
- [ ] Hostname, exit node, routes, and MagicDNS settings round-trip to the controller and back.
- [ ] MagicDNS is on by default for a fresh install.
- [ ] Login screen shows the auth URL/QR and waits for control-plane approval, then reflects the connected state.
- [ ] Logout offers the wipe/keep-keys choice.
- [ ] Errors appear inline and key events trigger a notification.

## Out of scope

- Tailnet-wide ACL/policy management (device-local controls only).
- Server roles (subnet router / exit node), ACL tags, and Tailscale SSH (deferred).