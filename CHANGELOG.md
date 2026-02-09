# Changelog

## v0.11.1

### Fixes
- **Fix plugin manifest compatibility for Helm v3 and v4** — ship `plugin.yaml` in Helm 3 (legacy) format so both versions can parse it; the install hook swaps to the Helm 4 manifest when needed (#234, #235, #237)

### Tests
- Add plugin manifest compatibility tests that validate YAML against both Helm v3 and v4 parsing logic, preventing future regressions

## v0.11.0

### Breaking Changes
- **Helm v2 support removed** — all Helm v2 code and references have been dropped

### Features
- **Helm v4 support** — dual Helm v3/v4 API support using wrapper pattern with runtime version detection (#226)
- **arm64/aarch64 architecture support** added to builds and install script
- Separate plugin manifests for Helm v3 and v4, with install hook auto-detection

### Dependencies
- Go 1.25
- Helm v3.18.5 (#230)
- Helm v4.0.4
- cobra 1.8.1
- protobuf 1.33.0
- Removed unused `docker/docker` indirect dependency

### Fixes
- CVE-2023-2253 fixed in golang deps (#197)
- GoReleaser deprecation warning resolved

### Other
- README refactored for clearer install instructions (#214)
- Acceptance tests updated for Helm v3 and v4
