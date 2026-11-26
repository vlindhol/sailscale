# Specs

Forward-looking design documents for Sailscale. Each spec is a self-contained unit of work: technical/architectural decisions and business rules for one feature, implementable in a few (preferably one) PRs.

## Conventions

Spec filenames follow `xxx-<state>-<title>.md`:
- `xxx` — an incrementing number.
- `<state>` — one of `wip`, `done`, `deprecated`.
- A `wip` spec is being started and should not stay in that state long (a PR or two).

Reading the `done` specs in order yields the project's decision history. If two specs conflict, the higher-numbered one wins. A spec that is mostly overridden by later specs moves to `deprecated`.

## Current specs

| # | Title | State |
|---|-------|-------|
| 001 | Privilege model (run without root) | `wip` |
| 002 | Core connectivity and engine control (baseline) | `wip` |
| 003 | ConnMan VPN plugin (supervisor + status relay) | `wip` |
| 004 | Authentication flow (interactive login) | `wip` |
| 005 | Standalone app UI (status, peers, settings) | `wip` |
| 006 | Settings-app integration | `wip` |

## Not in specs/

- **Reference about the environment** (ConnMan internals, Tailscale privilege modes, D-Bus API) and **current state / business rules** live in `corpus/`, per the repo layout in the root `AGENTS.md`.
- The formerly separate `decisions.md` log has been dissolved: accepted decisions are encoded as business rules in the relevant specs/corpus; history lives in git.