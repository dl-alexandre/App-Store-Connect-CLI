# ASC - App Store Connect CLI

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-yellow?style=for-the-badge" alt="License">
  <img src="https://img.shields.io/badge/Homebrew-compatible-blue?style=for-the-badge" alt="Homebrew">
</p>

A **fast**, **lightweight**, and **AI-agent friendly** CLI for App Store Connect. Ship iOS apps with zero friction.

## Why ASC?

| Problem | Solution |
|---------|----------|
| Manual App Store Connect work | Automate everything from CLI |
| Slow, heavy tooling | Go binary, fast startup |
| Not AI-agent friendly | JSON output, explicit flags, clean exit codes |

## Quick Start

### Install

```bash
# Via Homebrew (coming soon)
brew install rudrank/tap/asc

# Install script (macOS/Linux)
curl -fsSL https://raw.githubusercontent.com/rudrankriyam/App-Store-Connect-CLI/main/install.sh | bash

# Installs to ~/.local/bin by default (ensure it's on your PATH)

# Or build from source
git clone https://github.com/rudrankriyam/App-Store-Connect-CLI.git
cd App-Store-Connect-CLI
make build
./asc --help
```

### Authenticate

```bash
# Register your App Store Connect API key
asc auth login \
  --name "MyApp" \
  --key-id "ABC123" \
  --issuer-id "DEF456" \
  --private-key /path/to/AuthKey.p8
```

Generate API keys at: https://appstoreconnect.apple.com/access/integrations/api

Credentials are stored in the system keychain when available, with a local config fallback
at `~/.asc/config.json` (restricted permissions).
Environment variable fallback:
- `ASC_KEY_ID`
- `ASC_ISSUER_ID`
- `ASC_PRIVATE_KEY_PATH`

App ID fallback:
- `ASC_APP_ID`

## Commands

### Agent Quickstart

- Prefer `--json` (default) for machine parsing; add `--pretty` when debugging.
- Use `--limit` + `--next "<links.next>"` for pagination across all list commands.
- Sort with `--sort` (prefix `-` for descending):
  - Feedback/Crashes: `createdDate` / `-createdDate`
  - Reviews: `rating` / `-rating`, `createdDate` / `-createdDate`
  - Apps: `name` / `-name`, `bundleId` / `-bundleId`
  - Builds: `uploadedDate` / `-uploadedDate`

### TestFlight

```bash
# List beta feedback (JSON - best for AI agents)
asc feedback --app "123456789" --json

# Filter feedback by device model and OS version
asc feedback --app "123456789" --device-model "iPhone15,3" --os-version "17.2" --json

# Filter feedback by platform/build/tester
asc feedback --app "123456789" --app-platform IOS --device-platform IOS --build "BUILD_ID" --tester "TESTER_ID" --json

# Get crash reports (table format - for humans)
asc crashes --app "123456789" --output table

# Get crash reports (markdown - for docs)
asc crashes --app "123456789" --output markdown

# Limit results per page (pagination)
asc crashes --app "123456789" --limit 25 --json

# Sort crashes by created date (newest first)
asc crashes --app "123456789" --sort -createdDate --limit 5 --json

# Fetch next page
asc crashes --next "<links.next>" --json
```

### App Store

```bash
# List customer reviews (JSON - best for AI agents)
asc reviews --app "123456789" --json

# Filter by stars (table format - for humans)
asc reviews --app "123456789" --stars 1 --output table

# Filter by territory (markdown - for docs)
asc reviews --app "123456789" --territory US --output markdown

# Sort reviews by created date (newest first)
asc reviews --app "123456789" --sort -createdDate --limit 5 --json

# Fetch next page using links.next
asc reviews --next "<links.next>" --json
```

### Apps & Builds

```bash
# List apps (useful for finding app IDs)
asc apps --json

# Sort apps by name or bundle ID
asc apps --sort name --json
asc apps --sort -bundleId --json

# List builds for an app
asc builds --app "123456789" --json

# Sort builds by upload date (newest first)
asc builds --app "123456789" --sort -uploadedDate --json

# Fetch next page
asc apps --next "<links.next>" --json
asc builds --next "<links.next>" --json
```

### Utilities

```bash
# Print version information
asc version
asc --version
```

### Output Formats

| Format | Flag | Use Case |
|--------|------|----------|
| JSON (minified) | `--json` | AI agents, scripting |
| Table | `--output table` | Humans in terminal |
| Markdown | `--output markdown` | Humans, documentation |

### Authentication

```bash
# Check authentication status
asc auth status

# Logout
asc auth logout
```

## Design Philosophy

### Explicit Over Cryptic

```bash
# Good - self-documenting
asc reviews --app "MyApp" --stars 1 --json

# Avoid - cryptic flags (hypothetical, not supported)
# asc reviews -a "MyApp" -s 1
```

### AI-Agent Friendly

All commands output minified JSON by default for easy parsing by AI agents:

```bash
asc feedback --app "123456789" --json | jq '.data[].attributes.comment'
```

The `--json` flag outputs minified JSON (one line per response) to minimize token usage. Use `--output table` or `--output markdown` for human-readable output.

### No Interactive Prompts

Everything is flag-based for automation:

```bash
# Non-interactive (good for CI/CD and AI)
asc feedback --app "123456789" --json

# No prompts, no waiting
```

## Installation

### Homebrew (macOS)

```bash
# Add tap
brew tap rudrank/tap/asc

# Install
brew install asc
```

### From Source

```bash
git clone https://github.com/rudrankriyam/App-Store-Connect-CLI.git
cd App-Store-Connect-CLI
make build
make install  # Installs to /usr/local/bin
```

## Documentation

- [CLAUDE.md](CLAUDE.md) - Development guidelines for AI assistants
- [PLAN.md](PLAN.md) - Detailed roadmap and feature list
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contribution guidelines

## Roadmap

| Version | Features |
|---------|----------|
| v0.1 | Feedback, crashes, reviews |
| v0.2 | Apps, builds management |
| v0.3 | Beta testers, groups |
| v0.4 | Localizations |
| v0.5 | App submission |
| v1.0 | Full feature set |

See [PLAN.md](PLAN.md) for detailed roadmap.

## Security

- Credentials stored in the system keychain when available
- Local config fallback with restricted permissions
- Private key content never stored, only path reference
- Environment variables as fallback

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Author

[Rudrank Riyam](https://github.com/rudrankriyam)

---

<p align="center">
  Built with Go and Claude Code
</p>
