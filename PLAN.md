# ASC CLI - Project Plan

## Vision

A fast, AI-agent-friendly CLI for App Store Connect that enables developers to ship iOS apps with zero friction.

## Current Reality (v0.1 - Implemented, Validated Locally)

**Last Updated:** 2026-01-23

### What Doesn't Work Yet

- Integration tests are opt-in and require real credentials
- End-to-end build uploads are not automated (uploading IPA parts and committing the upload)
- Live coverage for mutating endpoints is incomplete (submit/build upload, beta group/tester mutations, localization create/update/delete)

## Roadmap (Remaining Work)

### Build Uploads: End-to-End

- Upload IPA parts to presigned URLs returned by build upload file creation
- Commit the upload with checksums
- Optional `--wait` to poll upload/build processing status

### Live Coverage

- Add opt-in integration tests for mutating endpoints:
  - submit/build upload
  - beta group and tester mutations
  - localization upload/create/update/delete

### Future Enhancements

- Interactive mode
- Plugins
- AI summarization
- Auto-responder to reviews
- Multi-account support
- Optional web UI

## Success Metrics

- Install via Homebrew: `brew install rudrank/tap/asc`
- Average startup time: < 50ms
- JSON output is default; use `--output` for table/markdown
- 80%+ test coverage
- Zero security vulnerabilities
