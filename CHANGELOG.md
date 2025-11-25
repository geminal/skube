# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.2] - 2025-11-25

### Added
- **Version Command**: Added `skube --version`, `-v`, and `version` commands to display the current version.
- **Help Text Updates**: Updated `skube help` to list all supported resources (nodes, configmaps, secrets, ingresses, pvcs).

### Fixed
- Fixed fmt.Sprintf argument count mismatch in help text causing build failure

### Documentation
- **README**: Complete overhaul with all commands, natural language features, and dynamic completion details.

## [0.2.0] - 2025-11-25

### Added
- **Natural Language Support**: Talk to Kubernetes naturally with 30+ command synonyms
  - Stop words filtering (`the`, `a`, `please`, `me`, `for`, etc.)
  - Command synonyms: `list`/`show`/`give` → `get`, `tail`/`monitor` → `logs`, etc.
  - Smart context inference for implicit namespaces and resources
- **Dynamic Shell Completion**: Autocomplete with real cluster resources
  - Queries actual namespaces, pods, deployments, and services from your cluster
  - Context-aware suggestions based on command structure
  - No more hardcoded namespace suggestions
- **Security Enhancements**:
  - kubectl validation on startup with OS-specific installation instructions
  - Enhanced input sanitization blocking shell injection characters
  - Flag injection prevention
  - Length limiting to prevent DoS
- **New Kubernetes Resources**:
  - Nodes (`nodes`, `no`)
  - ConfigMaps (`configmaps`, `cm`)
  - Secrets (`secrets`)
  - Ingress (`ingresses`, `ing`)
  - PersistentVolumeClaims (`persistentvolumeclaims`, `pvc`)
- **Documentation**:
  - CONTRIBUTING.md with comprehensive contributor guidelines
  - Enhanced README with smart autocomplete section
  - Updated COMMAND_REFERENCE.md with new resources
- **Testing**:
  - Comprehensive test suite with 78.9% parser coverage
  - Security sanitization tests
  - Natural language parsing tests
  - Resource parsing tests

### Changed
- **Parser Refactoring**: Complete rewrite using map-based lookups
  - Reduced code duplication by ~40%
  - Modular helper functions for better maintainability
  - Easier to extend with new commands and resources
- **Improved Error Messages**: kubectl not found now shows OS-specific installation instructions

### Fixed
- `.gitignore` now correctly allows `cmd/skube/` directory
- Parser now handles preposition-less commands (`skube pods qa` works)
- Type correction for natural language (e.g., "restart backend deployment")

## [0.1.0] - 2024-XX-XX

### Added
- Initial release
- Basic natural language parsing for kubectl commands
- Support for core resources: pods, deployments, services, namespaces
- Log streaming with follow and prefix options
- Shell access to pods
- Deployment operations: restart, scale, rollback
- Port forwarding
- Resource describe and delete
- Configuration management
- Shell completion for zsh and bash

[Unreleased]: https://github.com/geminal/skube/compare/v0.2.2...HEAD
[0.2.2]: https://github.com/geminal/skube/compare/v0.2.0...v0.2.2
[0.2.0]: https://github.com/geminal/skube/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/geminal/skube/releases/tag/v0.1.0
