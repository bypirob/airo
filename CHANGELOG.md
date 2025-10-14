# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added

- Direct image copy deployment mode via `transport: copy` in `env.yml`.
- `--config` flag to choose config directory for `env.yml`, `compose.yml`, and `Caddyfile`.
- `version` subcommand to print version, commit, and build date.
- Embedded build metadata via Makefile ldflags (version, commit, build date).
- Default config values: `transport: registry` and `user: airo` when omitted.
