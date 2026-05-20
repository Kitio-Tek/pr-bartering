# Contributing

Thanks for your interest in pr-bartering. This document describes how to propose
changes.

## Getting started

1. Read [DEVELOPER.md](DEVELOPER.md) and get `make check` passing locally.
2. Open an issue for anything more than a small fix, so the approach can be
   discussed before you spend time on it.

## Branching and commits

- Branch from `main`. Keep one logical change per branch.
- Write commit messages in the [Conventional Commits](https://www.conventionalcommits.org)
  style, for example `fix: re-replicate data when a proof times out`. Common
  prefixes used here are `feat`, `fix`, `refactor`, `test`, `docs`, `build`,
  `ci`, `chore`, and `security`.
- Keep the subject line short and in the imperative mood, and use the body to
  explain what changed and why.

## Sign-off

All commits must be signed off under the
[Developer Certificate of Origin](https://developercertificate.org). Add the
sign-off line with `git commit -s`:

```
Signed-off-by: Your Name <you@example.com>
```

## Pull requests

- Fill in the pull-request template: summary, motivation, the testing you did,
  and the linked issue.
- A pull request must pass CI: build, race tests, lint, gosec, and govulncheck.
- Do not weaken or disable a check to make the pipeline pass. If a security
  finding is a genuine false positive, suppress it inline with a `#nosec` comment
  that explains why.
- Add tests with your change. New behaviour needs unit tests, including the
  relevant edge cases (empty input, missing peers, malformed messages, timeouts).
- Update the documentation when behaviour or configuration changes. Every claim
  in the README and `docs/` should be backed by code that works.
- If the change is backwards-incompatible, say so in the pull request and label
  it `breaking-change` so it is captured in the release notes.

## Review

A maintainer reviews each pull request. Expect questions about edge cases and
test coverage. Keep the discussion on the pull request so the rationale stays
with the code.

## Reporting bugs and requesting features

Use the issue templates. For security problems, follow [SECURITY.md](SECURITY.md)
instead of opening a public issue.
