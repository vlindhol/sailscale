# Spec 004 — Authentication flow (interactive login)

## Goal

Let a user log in and authenticate the device to a tailnet through an interactive flow, without any persistent secret material handled by Sailscale itself.

## Context

Tailscale device authentication is out-of-band (a browser approval), not a username/password prompt. The device initiates a login, shows the user a URL/QR code, and the user approves in their browser; the control plane then marks the device authenticated. The ConnMan `Agent` credential mechanism is therefore **not** used.

## Decisions (technical)

- **Interactive login (URL + QR) is the only mechanism for v1.** Auth-key entry is out of scope.
- The short-lived `sailscale-ctl` controller asks `tailscaled` (via LocalAPI) to start interactive auth and returns the auth URL; the UI renders it as both a URL and a QR code with a "waiting for login" state.

## Business rules

- The UI shows the login URL and QR, and transitions from "waiting for login" to "connected" once the control-plane approval completes.
- On success, the profile is persisted so the device rejoins on next connect.
- On auth timeout or failure, the error is surfaced inline with an option to retry.
- Logging out signs the device out of the tailnet, with an optional full wipe of local state (a "keep keys" checkbox defaulting to wipe for quick re-auth).

## Acceptance criteria

- [ ] Login produces a working auth URL and QR code.
- [ ] Approving in a browser results in the device joining the tailnet (100.x address assigned).
- [ ] "Waiting for login" resolves to "connected" without a manual refresh.
- [ ] Timeout/failure surfaces inline with a retry path.
- [ ] Logout signs out; with the wipe option it removes local state/keys.

## Out of scope

- Auth-key (pre-auth/OWA-key) entry.
- Tailnet-wide identity/policy management (that lives in the admin console).