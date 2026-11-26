# Spec 002 — Core connectivity and engine control (baseline)

## Goal

Establish the foundation: run the official Tailscale engine (`tailscaled`) on Sailfish OS and control it from a short-lived Go helper, so that a device can join a tailnet, obtain a 100.64.0.0/10 address, reach peers, and use advertised subnet routes.

## Context

Sailscale wraps the official Tailscale engine rather than reimplementing it. `tailscaled` is the resident daemon that holds node state, talks to the Tailscale control plane (`login.tailscale.com`), and drives the tunnel. It is **not** to be reimplemented.

## Decisions (technical)

- **Engine lifecycle is supervised by a ConnMan VPN plugin**, not by a dedicated long-running Sailscale daemon. There is therefore **no resident Sailscale controller daemon**; `tailscaled` is spawned and stopped by the ConnMan plugin.
- **Control happens via a short-lived Go binary** (`sailscale-ctl`) that runs on demand and exits. It is a thin backend-for-frontend over the bundled `tailscale` CLI: it execs only the CLI's local features (`up`/`down`/`status`/`version`) and re-emits a stable JSON contract, isolating all Tailscale-specific knowledge to that single translation layer. It is not a service.
- **The engine is bundled** in the Sailscale RPM (a pinned `tailscaled`/`tailscale` build). No separate Sailfish-tailscale package dependency.
- **Kernel networking is the default.** The engine runs in kernel-tun mode (full performance, all protocols), owned by the ConnMan plugin's helper slot. Userspace-networking is *not* the default and is used only as a degrade-gracefully fallback.

## Business rules

- The user logs in and authenticates the device to a tailnet. Interactive login (URL + QR) is the primary mechanism; auth-key entry is out of scope for v1.
- After login the device obtains its 100.64.0.0/10 CGNAT address.
- The device can reach other tailnet nodes over the tunnel, and be reached by them, subject to tailnet ACLs.
- The device can accept and use peer- and subnet-router-advertised routes (client-side route acceptance).
- The user can connect and disconnect the VPN from the system UI.
- Supported platform is Sailfish OS 5.x (Jolla C2 and later). Sailfish OS 4.x is out of scope for v1.
- The app surfaces the bundled engine version in the UI.

## Acceptance criteria

- [ ] `tailscaled` starts and stays running under supervision, as a non-root user.
- [ ] `sailscale-ctl` can drive `tailscale up/down/status/version` by exec'ing the bundled CLI's local features and emitting a stable JSON contract.
- [ ] After login the device appears in the tailnet with a 100.x address.
- [ ] Peers are reachable and can reach this device (ACL permitting).
- [ ] Advertised routes are accepted and traffic flows over them.
- [ ] Connect/disconnect reflects correctly in the system connectivity UI.
- [ ] Version is discoverable and surfaced.

## Out of scope

- Server roles (subnet router / exit node), ACL tags at join, and Tailscale SSH are deferred to later specs.