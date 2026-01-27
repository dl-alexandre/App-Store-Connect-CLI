# OpenAPI snapshot (offline)

This folder keeps an offline snapshot of the App Store Connect OpenAPI spec for
agents that cannot access the internet.

## Files

- `latest.json`: full OpenAPI spec snapshot (see source below)
- `paths.txt`: generated path+method index for quick existence checks

## Source

The snapshot comes from the community-maintained OpenAPI repo:
`https://github.com/EvanBacon/App-Store-Connect-OpenAPI-Spec`

## Update process

1. Replace `latest.json` with a newer spec file.
2. Run `scripts/update-openapi-index.py` to regenerate `paths.txt`.
3. Update the "Last synced" date below.

Last synced: 2026-01-27
