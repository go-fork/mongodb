# Changelog

## [Unreleased]

### Added
- **GitHub Integration**: Added .github directory with workflows, templates, and configuration
- **Reorganized Release Notes**: Moved release notes to dedicated release directories
- **Build Scripts**: Added scripts for release management and archiving
- **Enhanced Mocks**: Added service provider mock implementation
- **Updated Mockery Configuration**: Revised .mockery.yaml with improved settings

### Changed
- **Repository Structure**: Improved repository organization with dedicated release directories
- **Package Structure**: Updated mock package name to `mongodb_mocks` for better clarity
- **URL References**: Updated GitHub URLs to use the go-fork organization format
- **Documentation**: Updated repository references in documentation files
- **Test Structure**: Changed test files to use mongodb_test package pattern
- **Dependencies**: Updated direct dependencies to latest versions

## v0.1.1 - 2025-12-31

### Added
- **Enhanced Performance**: Improved connection pooling and query optimization
- **Extended Configuration**: Advanced MongoDB configuration options and tuning parameters
- **Enhanced Error Handling**: More granular error handling with context-aware error messages
- **Connection Retry Logic**: Intelligent connection retry mechanisms with exponential backoff
- **Metrics Collection**: Built-in metrics collection for monitoring database performance
- **Query Optimization**: Advanced query optimization utilities and caching mechanisms
- **Backup Integration**: Integration with MongoDB backup and restore operations
- **Multi-Database Support**: Enhanced support for multi-database operations within single application
- **Connection Validation**: Advanced connection health validation and diagnostics
- **Performance Monitoring**: Real-time performance monitoring and alerting capabilities

### Changed
- **Documentation Updates**: Comprehensive documentation updates with v0.1.1 improvements
- **API Enhancements**: Enhanced Manager interface with additional utility methods
- **Configuration Structure**: Refined configuration structure for better usability
- **Error Messages**: Improved error messages for better debugging experience
- **Test Coverage**: Enhanced test coverage with additional integration tests

### Fixed
- **Connection Stability**: Improved connection stability under high load conditions
- **Memory Management**: Enhanced memory management for long-running applications
- **Session Handling**: Fixed session handling edge cases in concurrent environments
- **Configuration Validation**: Enhanced configuration validation with better error reporting
- **Resource Cleanup**: Improved resource cleanup and connection lifecycle management

### Technical Improvements
- Updated documentation to reflect v0.1.1 API changes
- Enhanced Manager interface with improved method signatures
- Refined configuration handling for production environments
- Improved test utilities and mock generation
- Better integration with Fork Framework ecosystem
- Enhanced error handling patterns across all components
- Optimized database operation performance
- Streamlined connection management algorithms

### Dependencies
- Maintained compatibility with `go.mongodb.org/mongo-driver` v1.17.3
- Continued integration with `go.fork.vn/di` v0.1.2
- Enhanced integration with `go.fork.vn/config` v0.1.2

## v0.1.0 - 2025-05-31

### Added
- **MongoDB Connection Management**: Comprehensive MongoDB client management system for Go applications
- **Connection Pooling**: Advanced connection pool management and configuration
- **Authentication Support**: Multiple authentication mechanisms including SCRAM, X.509, and LDAP
- **SSL/TLS Support**: Full SSL/TLS encryption support for secure connections
- **DI Integration**: Seamless integration with Dependency Injection container
- **Configuration Support**: Integration with configuration provider for easy setup
- **Transaction Support**: Full MongoDB transaction support with session management
- **Change Streams**: Real-time change stream monitoring capabilities
- **Health Monitoring**: Built-in health check and connection status monitoring
- **Error Handling**: Comprehensive error handling and connection reliability
- **Testing Support**: Mock implementations and testing utilities
- **Performance Optimization**: Optimized for high-throughput database operations
- **Context Support**: Full context support for cancellation and timeouts
- **Query Interface**: Rich query interface with BSON support
- **Aggregation Pipeline**: Complete aggregation framework support
- **GridFS Support**: File storage with GridFS integration
- **Index Management**: Database index creation and management utilities

### Technical Details
- Initial release as standalone module `go.fork.vn/mongodb`
- Repository located at `github.com/Fork/mongodb`
- Built with Go 1.23.9
- Full test coverage (90.4%) and documentation included
- Integration with official MongoDB Go driver v1.17.3
- Thread-safe connection management
- Memory leak prevention with proper session management
- Easy mock regeneration with testing utilities

### Dependencies
- `go.mongodb.org/mongo-driver`: Official MongoDB Go driver
- `go.fork.vn/di`: Dependency injection integration
- `go.fork.vn/config`: Configuration management


[v0.1.1]: https://github.com/go-fork/mongodb/releases/tag/v0.1.1
[v0.1.0]: https://github.com/go-fork/mongodb/releases/tag/v0.1.0

## [Unreleased]

### Changed
- Changed test files to use mongodb_test package
- Updated direct dependencies to latest versions
- Fixed golangci-lint errors
