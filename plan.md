# Coverage Expansion Plan (audit/fix-security-and-test-gaps)

Goal: eliminate current `0.0%` coverage packages by adding package-local tests in phases, then re-run short/full suites.

## Baseline (captured)

- [x] Run `go test -short -coverprofile=coverage.out ./...`
- [x] Capture current zero-coverage package list
- [x] Re-check baseline list after each phase

## Phase 0 - Planning and Tracking

- [x] Create `plan.md`
- [ ] Keep this file updated after every phase commit

## Phase 1 - Core Utility and Small Command Packages

- [x] `internal/asc/types`
- [x] `internal/cli/shared/suggest`
- [x] `internal/cli/completion`
- [x] `internal/cli/registry`
- [x] `internal/cli/crashes`
- [x] `internal/cli/feedback`
- [x] `internal/cli/submit`
- [x] `internal/cli/routingcoverage`
- [x] Commit Phase 1

## Phase 2 - Identity/Metadata Command Families

- [x] `internal/cli/accessibility`
- [x] `internal/cli/actors`
- [x] `internal/cli/agerating`
- [x] `internal/cli/agreements`
- [x] `internal/cli/categories`
- [x] `internal/cli/eula`
- [x] `internal/cli/nominations`
- [x] `internal/cli/merchantids`
- [x] `internal/cli/passtypeids`
- [x] Commit Phase 2

## Phase 3 - App Distribution/Release Command Families

- [x] `internal/cli/androidiosmapping`
- [x] `internal/cli/app_events`
- [x] `internal/cli/appclips`
- [x] `internal/cli/assets`
- [x] `internal/cli/betaapplocalizations`
- [x] `internal/cli/betabuildlocalizations`
- [x] `internal/cli/buildbundles`
- [x] `internal/cli/buildlocalizations`
- [x] `internal/cli/localizations`
- [x] `internal/cli/prerelease`
- [x] `internal/cli/productpages`
- [x] Commit Phase 3

## Phase 4 - Commerce/Security/Operational Command Families

- [x] `internal/cli/encryption`
- [x] `internal/cli/offercodes`
- [x] `internal/cli/performance`
- [x] `internal/cli/promotedpurchases`
- [x] `internal/cli/reviews`
- [x] `internal/cli/winbackoffers`
- [x] `internal/cli/notarization`
- [x] Commit Phase 4

## Phase 5 - Large Surface Command Families

- [x] `internal/cli/gamecenter`
- [x] `internal/cli/xcodecloud`
- [x] Commit Phase 5

## Phase 6 - Remaining Non-CLI Root Packages

- [ ] `github.com/rudrankriyam/App-Store-Connect-CLI` (main package)
- [ ] Commit Phase 6

## Validation Gate (after each phase and at end)

- [ ] `go test -short ./...`
- [ ] `go test -short -coverprofile=coverage.out ./...`
- [ ] verify previously-targeted packages are no longer `0.0%`
- [ ] `make test`
- [ ] `make lint`
