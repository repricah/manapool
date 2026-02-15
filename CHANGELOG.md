# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.0] - 2026-02-15

### Changed

- **BREAKING (Configuration)**: Updated inventory API limits to match Manapool API v0.24.1 specification
  - Maximum inventory limit increased from 500 to 10,000 items per request
  - Default inventory limit changed from 500 to 100 items per request
  - `IterateInventory()` helper now uses 1,000 items per page (was 500) for better performance
  
### Migration Notes

This change is **backward compatible** for most users:
- ✅ Code using limits ≤ 500 continues to work without changes
- ✅ Code using default limit (not specifying `Limit: 0`) will automatically benefit from the new default
- ⚠️ Code relying on the **implicit default of 500** will now get 100 items per page instead

See [MIGRATION_GUIDE.md](./MIGRATION_GUIDE.md) for detailed upgrade instructions.

### Fixed

- Aligned client validation with current API specification limits

## [0.2.0] - 2025-XX-XX

### Added

- Initial pre-release version
- Seller inventory endpoints
- Type-safe API definitions
- Automatic rate limiting
- Automatic retries with exponential backoff
- Context support for cancellation and timeouts
- Comprehensive error handling with specific error types
- 96.5% test coverage

### Features

- Get seller account information
- List seller inventory with pagination
- Lookup inventory items by TCGPlayer SKU
- Helper function for iterating all inventory items

[Unreleased]: https://github.com/repricah/manapool/compare/v0.3.0...HEAD
[0.3.0]: https://github.com/repricah/manapool/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/repricah/manapool/releases/tag/v0.2.0
