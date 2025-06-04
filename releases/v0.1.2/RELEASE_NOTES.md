# Release Notes - v0.1.2

## Overview
This release focuses on GitHub integration, repository structure improvements, and enhanced development experience with comprehensive CI/CD workflows and automation.

## What's New
### 🚀 Features
- **GitHub Integration**: Complete .github directory setup with CI/CD workflows, issue templates, and automation
- **Release Management**: Automated release workflows with proper versioning and changelog management
- **Dependency Updates**: Automated dependency management with Dependabot integration
- **Enhanced Mocks**: New service provider mock implementation with improved generation
- **Build Scripts**: Comprehensive release management and archiving scripts

### 🔧 Improvements
- **Repository Structure**: Organized release notes into dedicated version directories
- **Package Structure**: Updated mock package naming to `mongodb_mocks` for better clarity
- **URL References**: Updated all GitHub URLs to use go-fork organization format
- **Test Structure**: Improved test organization with mongodb_test package pattern
- **Mock Generation**: Enhanced .mockery.yaml configuration for better mock generation

### 📚 Documentation
- **Repository References**: Updated all documentation to reflect go-fork organization
- **Release Documentation**: Structured release notes in dedicated directories
- **Scripts Documentation**: Added comprehensive documentation for build and release scripts

## Breaking Changes
### ⚠️ Important Notes
- **Package Names**: Mock package renamed from previous naming to `mongodb_mocks`
- **URL Changes**: Repository URLs changed from previous organization to go-fork

## Migration Guide
See [MIGRATION.md](./MIGRATION.md) for detailed migration instructions.

## Dependencies
### Updated
- `go.fork.vn/config`: v0.1.2 → v0.1.3
- `go.fork.vn/di`: v0.1.2 → v0.1.3

### Maintained
- `go.mongodb.org/mongo-driver`: v1.17.3
- `github.com/stretchr/testify`: v1.10.0

## Performance
- Benchmark improvement: X% faster in scenario Y
- Memory usage: X% reduction in scenario Z

## Security
- Security fix for vulnerability X
- Updated dependencies with security patches

## Testing
- Added X new test cases
- Improved test coverage to X%

## Contributors
Thanks to all contributors who made this release possible:
- @contributor1
- @contributor2

## Download
- Source code: [go.fork.vn/mongodb@v0.1.2]
- Documentation: [pkg.go.dev/go.fork.vn/mongodb@v0.1.2]

---
Release Date: 2025-06-04
