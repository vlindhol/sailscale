# Tailscale bundle

Sailscale does not build Tailscale; it bundles a pinned build in its RPM.

## Pinned build

- Track: `stable`, version `1.104.7`, architecture `arm64`.
- Source: `https://pkgs.tailscale.com/stable/tailscale_1.104.7_arm64.tgz`
  (SHA-256 `d6bed2720310b328b861f943d7b4016c5358a9bad60575eb51e8143b79068484`).
- Contains `tailscaled` (the resident daemon) and `tailscale` (the CLI, which
  `sailscale-ctl` execs for its local features).

## Rules

- Tailscale is bundled in Sailscale — there is **no** separate
  Sailfish-tailscale package dependency.
- The pinned version lives in `tailscale/VERSION`, the single source of truth;
  the build injects it into `sailscale-ctl` so the bundled version can be
  surfaced.
- The bundled version is surfaced to the user (see `sailscale-ctl version`).
- Kernel networking (tun) is the default; userspace networking is only a
  degrade-gracefully fallback.