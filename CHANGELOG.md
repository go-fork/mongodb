# Changelog

## [Unreleased]

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
- Repository located at `github.com/go-fork/mongodb`
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

[Unreleased]: https://github.com/go-fork/mongodb/compare/v0.1.0...HEAD
[v0.1.0]: https://github.com/go-fork/mongodb/releases/tag/v0.1.0
