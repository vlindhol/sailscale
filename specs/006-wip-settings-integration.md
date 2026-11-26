# Spec 006 — Settings-app integration

## Goal

Expose Sailscale status and controls inside the Sailfish OS general Settings app, in addition to the standalone app, so users can connect/disconnect and see state without opening the full app.

## Context

Sailfish OS allows third-party integration into the Settings app via plugins such as `nemo-qml-plugin-systemsettings`. This is a lightweight companion to the standalone app; the VPN connection itself already appears in the system connectivity UI because Sailscale registers a ConnMan VPN provider.

## Decisions (technical)

- Add a Settings-app entry for Sailscale that folds status/settings into the general Settings app.
- The Settings surface is **lightweight**: status (online/IP/peers), Connect/Disconnect, and a deep-link to open the standalone app for advanced configuration.

## Business rules

- The Settings entry shows the same authoritative status as the standalone app.
- Connect/Disconnect from Settings maps to the same ConnMan/engine path as everywhere else.
- Advanced configuration is intentionally not duplicated in Settings; it opens the standalone app.

## Acceptance criteria

- [ ] A Sailscale entry appears in the Settings app.
- [ ] Status shown in Settings matches the standalone app and the connectivity UI.
- [ ] Connect/Disconnect works from Settings.
- [ ] The "open app" deep-link launches the standalone app.

## Out of scope

- Recreating the full settings UI inside the Settings app.