# Security policy

## Supported versions

pr-bartering is research software under active development. Security fixes are
applied to the `main` branch and the most recent tagged release.

## Reporting a vulnerability

Please report vulnerabilities privately. Do not open a public issue, pull
request, or discussion for a security problem.

Use GitHub's private reporting:
[open a security advisory](https://github.com/Kitio-Tek/pr-bartering/security/advisories/new).

Include as much of the following as you can:

- A description of the issue and its impact.
- The affected component or package.
- Steps to reproduce, or a proof of concept.
- Any suggested remediation.

## What to expect

- An acknowledgement of your report within five working days.
- An assessment of the issue and, if it is confirmed, a plan and timeline for a
  fix.
- Credit in the release notes for the fix, unless you ask to remain anonymous.

## Scope

This project shells out to the local `ipfs` binary and exchanges messages with
peers over TCP. Reports about the handling of untrusted peer input, the
proof-of-storage mechanism, and the bootstrap service are especially welcome.
