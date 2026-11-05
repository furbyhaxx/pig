<!--
SPDX-License-Identifier: MIT
-->

# Security Policy

This repository is a personal fork of [PiG](https://github.com/MichaelKinsy/PiG). Security reports for this fork go to its owner, not to PiG's maintainers.

## Supported versions

The latest published fork release receives security fixes. `main` is development code. Older releases can require an upgrade to receive a fix.

## Security boundary

The agent runs with the operating-system permissions of the user who starts it. It does not sandbox model output, shell commands, tools, extensions, skills, hooks, project instructions, or local configuration.

Treat the user account, writable files, environment, shell configuration, agent configuration and trusted project files as one local trust boundary. Use a container, virtual machine or another operating-system boundary when you need isolation.

A report must show how the fork grants access or crosses a security boundary. Expected behavior from already trusted local input is not a vulnerability by itself.

## Report a vulnerability

Do not report an undisclosed vulnerability in a public issue, pull request or discussion.

Use GitHub private vulnerability reporting on this repository. Include the affected version or commit, the operating system and installation method, the vulnerability type and impact, reproduction steps, known mitigations and any disclosure deadline. Do not include live credentials, customer data or unrelated confidential information.

## Out of scope

- prompt injection from content the user chose to trust;
- behavior of an untrusted extension, skill, hook, tool or project instruction that the user chose to load;
- local execution with the permissions of the user who started the agent;
- modification of files that the attacker could already write;
- malicious model output;
- intentionally weakened configuration;
- public exposure of a local agent process; and
- denial-of-service claims that require trusted local input and do not cross a resource boundary.

## Disclosure

The repository owner assesses reports and coordinates disclosure. No fixed response or remediation time is promised. A fix can require dependency updates, a patched release and refreshed release evidence.
