# Packaging

Sailscale is distributed as an RPM. `packaging/` records the inputs the package
is built from.

## Contents

- `sailscale.spec` — the RPM spec. It bundles the pinned Tailscale binaries and
  installs the `sailscale-ctl` control binary and the plugin `.conf`.
- `sailscale.conf` — the ConnMan VPN plugin configuration, installed to
  `/etc/connman/vpn-plugin/sailscale.conf`.

## Business rules

- The plugin `.conf` declares that the helper binary runs as the unprivileged
  `nemo` user with the `vpn` group and the `inet`/`net_admin` supplementary
  groups, so installing and connecting requires no root step.
- The Tailscale version in the spec (the `%define tailscale_version` macro) must
  agree with `tailscale/VERSION`.
- The bundled Tailscale binaries install as `tailscaled` under
  `libexec/sailscale/` and as `sailscale-tailscale` in the bindir (the CLI is
  renamed to avoid clashing with any system `tailscale`). `sailscale-ctl`
  execs this renamed CLI.