<!--
SPDX-License-Identifier: MIT
-->

# Contributing

This is a personal fork of [PiG](https://github.com/MichaelKinsy/PiG). It does not accept unsolicited external contributions, and it does not follow upstream PiG's contribution process.

## For the owner and agents

- Read the root [`AGENTS.md`](../AGENTS.md) and the `AGENTS.md` in the area you change.
- Keep fork decisions in `.project/specs/`; user-facing behavior changes update user docs and add a `changelog.d/` fragment.
- Do not introduce upstream parity requirements for fork behavior. See `.project/specs/0001-personal-distribution.md` D7–D10.
- Run the slice gates and `make check` before a release. Generated files are committed, never hand-edited.
- Commits follow the fork's Conventional Commits convention. Upstream PiG's DCO sign-off and signed-commit rules do not apply here.

Security reports follow [`.github/SECURITY.md`](SECURITY.md). Conduct reports follow [`.github/CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md).
