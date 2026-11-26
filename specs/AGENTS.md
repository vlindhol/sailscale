# How to write specs

* Specs are numbered documents of the form `xxx-<state>-<title>.md`. xxx is a number, `state` is other `wip`, `done`, `deprecated`.
* Each spec document is self-contained and does not reference other spec documents. It only contains technical/architectural decisions and business rules.
* One should be able to derive a list of decisions by reading the `done` spec files in order.
* `wip` spec files are most likely just started, a spec file shouldn't stay in that state for long, one or two PRs at the most.
* Spec files should be small enough to be easily digestible and implementable in a few (preferably one) PR, but big enough that they encompass entire features, not just small tweaks.
* If spec files are in conflict with each other, the higher numbered spec file wins.
* If a spec file is mostly overridden by other spec files, it moves from `done` to `deprecated`.
