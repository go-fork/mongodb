// Package MongoDB Provider v0.1.1 includes the following key features:
//   - Seamless integration with Go-Fork Framework's DI container
//   - Advanced connection management with connection pooling
//   - Enhanced authentication and SSL/TLS security
//   - Transaction support with automatic retry logic
//   - Change Streams for real-time data monitoring
//   - Built-in health checks and connection monitoring
//   - Comprehensive testing support with mockery integration
//   - Configuration-driven setup with environment variable support
//   - Performance optimizations and memory managementprovides MongoDB integration for Fork Framework v0.1.1.
//
// This package offers comprehensive MongoDB connectivity and management through
// the Fork Framework's dependency injection system, providing a clean and
// efficient way to work with MongoDB databases in Go applications.
//
// # Overview
//
// MongoDB Provider v0.1.1 includes the following key features:
//   - Seamless integration with Fork Framework's DI container
//   - Advanced connection management with connection pooling
//   - Multiple database connection support
//   - Enhanced authentication and SSL/TLS security
//   - Transaction support with automatic retry logic
//   - Change Streams for real-time data monitoring
//   - Built-in health checks and connection monitoring
//   - Comprehensive testing support with mockery integration
//   - Configuration-driven setup with environment variable support
//   - Performance optimizations and memory management
//
// # Architecture
//
// The MongoDB provider follows Fork Framework's service provider pattern:
//
//	Application
//	    ↓
//	ServiceProvider (mongodb.NewServiceProvider)
//	    ↓
//	Manager (mongodb.Manager)
//	    ├── Connection Pool Management
//	    ├── Database Operations
//	    └── Health Monitoring
//
// # Configuration
//
// MongoDB can be configured through YAML configuration files or programmatically.
// Example configuration in `configs/app.yaml`:
//
//	mongodb:
//	  uri: "mongodb://localhost:27017"
//	  database: "myapp"
//	  app_name: "my-golang-app"
//	  max_pool_size: 100
//	  min_pool_size: 5
//	  connect_timeout: 30000
//	  auth:
//	    username: "${MONGO_USERNAME}"
//	    password: "${MONGO_PASSWORD}"
//	    auth_source: "admin"
//	    auth_mechanism: "SCRAM-SHA-256"
//	  tls:
//	    enabled: true
//	    cert_file: "/path/to/cert.pem"
//	    key_file: "/path/to/key.pem"
//	    ca_file: "/path/to/ca.pem"
//
// # Service Provider Integration
//
// Register MongoDB provider with the Fork application:
//
//	import (
//	    "go.fork.vn/app"
//	    "go.fork.vn/config"
//	    "go.fork.vn/mongodb"
//	)
//
//	func main() {
//	    application := app.NewApplication()
//
//	    application.RegisterProviders(
//	        config.NewServiceProvider(),
//	        mongodb.NewServiceProvider(),
//	    )
//
//	    application.Boot()
//
//	    // Use MongoDB
//	    var manager mongodb.Manager
//	    application.Container().Make("mongodb.manager", &manager)
//	}
//
// # Manager Interface
//
// The Manager interface provides comprehensive MongoDB operations:
//
//	type Manager interface {
//	    Client() *mongo.Client
//	    Database() *mongo.Database
//	    DatabaseWithName(name string) *mongo.Database
//	    Collection(name string) *mongo.Collection
//	    CollectionWithDatabase(dbName, collectionName string) *mongo.Collection
//	    Config() *Config
//	    Ping(ctx context.Context) error
//	    Disconnect(ctx context.Context) error
//	    StartSession(opts ...*options.SessionOptions) (mongo.Session, error)
//	    UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error
//	    UseSessionWithTransaction(ctx context.Context, fn func(mongo.SessionContext) (interface{}, error), opts ...*options.TransactionOptions) (interface{}, error)
//	    HealthCheck(ctx context.Context) error
//	    Stats(ctx context.Context) (map[string]interface{}, error)
//	    ListCollections(ctx context.Context) ([]string, error)
//	    ListDatabases(ctx context.Context) ([]string, error)
//	    Watch(ctx context.Context, pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error)
//	    CreateIndex(ctx context.Context, collectionName string, model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error)
//	}
//
// # Basic Usage
//
//	// Get MongoDB manager
//	var manager mongodb.Manager
//	application.Container().Make("mongodb.manager", &manager)
//
//	// Work with collections
//	users := manager.Collection("users")
//
//	// Insert document
//	result, err := users.InsertOne(ctx, bson.M{
//	    "name": "John Doe",
//	    "email": "john@example.com",
//	})
//
//	// Find documents
//	cursor, err := users.Find(ctx, bson.M{"active": true})
//	defer cursor.Close(ctx)
//
// # Transaction Support
//
//	// Simple transaction
//	result, err := manager.UseSessionWithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
//	    users := manager.Collection("users")
//	    orders := manager.Collection("orders")
//
//	    // Perform multiple operations
//	    userResult, err := users.InsertOne(sc, user)
//	    if err != nil {
//	        return nil, err
//	    }
//
//	    orderResult, err := orders.InsertOne(sc, order)
//	    if err != nil {
//	        return nil, err
//	    }
//
//	    return map[string]interface{}{
//	        "user": userResult,
//	        "order": orderResult,
//	    }, nil
//	})
//
// # Change Streams
//
//	// Watch for changes
//	pipeline := mongo.Pipeline{
//	    bson.D{{"$match", bson.D{{"operationType", "insert"}}}},
//	}
//
//	stream, err := manager.WatchCollection(ctx, "users", pipeline)
//	if err != nil {
//	    return err
//	}
//	defer stream.Close(ctx)
//
//	for stream.Next(ctx) {
//	    var changeDoc bson.M
//	    if err := stream.Decode(&changeDoc); err != nil {
//	        continue
//	    }
//	    // Process change
//	}
//
// # Health Monitoring
//
//	// Check connection health
//	if err := manager.Ping(ctx); err != nil {
//	    log.Printf("MongoDB connection unhealthy: %v", err)
//	}
//
//	// Get connection statistics
//	stats, err := manager.Stats(ctx)
//	if err == nil {
//	    log.Printf("MongoDB stats: %+v", stats)
//	}
//
// # Index Management
//
//	// Create single field index
//	indexModel := mongo.IndexModel{
//	    Keys: bson.D{{"email", 1}},
//	    Options: options.Index().SetUnique(true),
//	}
//	indexName, err := manager.CreateIndex(ctx, "users", indexModel)
//
//	// List all indexes
//	cursor, err := manager.ListIndexes(ctx, "users")
//	defer cursor.Close(ctx)
//
// # Services Registered
//
// The ServiceProvider registers the following services in the DI container:
//   - "mongodb.manager" - Manager interface instance
//   - "mongodb" - Alias for manager
//   - "mongo.client" - Raw MongoDB client (*mongo.Client)
//   - "mongo" - Another alias for manager
//
// # Testing Support
//
// MongoDB Provider v0.1.1 includes comprehensive testing support with mockery:
//
//	func TestMongoOperations(t *testing.T) {
//	    mockManager := mocks.NewMockManager(t)
//
//	    mockManager.On("Ping", mock.Anything).Return(nil)
//	    mockManager.On("Collection", "users").Return(mockCollection)
//
//	    // Test your service
//	    service := NewUserService(mockManager)
//	    err := service.CreateUser(ctx, user)
//
//	    assert.NoError(t, err)
//	    mockManager.AssertExpectations(t)
//	}
//
// # Performance Optimization
//
// MongoDB Provider v0.1.1 includes several performance optimizations:
//   - Connection pooling with configurable pool sizes
//   - Automatic connection health monitoring
//   - Query optimization helpers
//   - Memory-efficient cursor handling
//   - Background connection cleanup
//
// # Security Features
//
//   - SSL/TLS encryption support
//   - Multiple authentication mechanisms
//   - Credential management with environment variables
//   - Network security configuration
//   - Access control integration
//
// # Dependencies
//
// MongoDB Provider v0.1.1 requires:
//   - Go 1.23.9+
//   - go.fork.vn/di v0.1.2
//   - go.fork.vn/config v0.1.2
//   - go.mongodb.org/mongo-driver v1.17.3+
//   - MongoDB Server 4.4+ (recommended 6.0+)
//
// # Migration from v0.1.0
//
// Key changes in v0.1.1:
//   - Enhanced Manager interface with comprehensive database operations
//   - Improved configuration structure with detailed options
//   - Updated dependencies (di v0.1.2, config v0.1.2)
//   - Better error handling and custom error types
//   - Performance optimizations and memory management
//   - Comprehensive testing support with mockery integration
//   - Enhanced index management capabilities
//   - Improved change streams and transaction support
//
// # Requirements
//
// - Go 1.23.9 or higher
// - MongoDB 4.4+ (recommended 6.0+)
// - Replica set configuration for transactions and change streams
//
// # Best Practices
//
//   - Always use context.Context with timeouts for MongoDB operations
//   - Configure appropriate connection pool sizes for your workload
//   - Use projections to limit returned data and improve performance
//   - Create proper indexes for frequently queried fields
//   - Monitor connection health and statistics regularly
//   - Use transactions for operations that require ACID guarantees
//   - Implement proper error handling for MongoDB-specific errors
//
// # Documentation
//
// Complete documentation is available in the docs/ directory:
//   - docs/index.md - Quick start guide
//   - docs/overview.md - Architecture and design
//   - docs/reference.md - Complete API reference
//   - docs/usage.md - Detailed usage examples
//
// For more information, visit: https://go.fork.vn/mongodb
package mongodb
