# Hướng Dẫn Sử Dụng MongoDB Provider v0.1.1

## Bắt Đầu Nhanh

### 1. Cài Đặt Package

```bash
go get go.fork.vn/mongodb@v0.1.1
```

### 2. Cấu Hình Cơ Bản

Tạo file `configs/app.yaml`:

```yaml
mongodb:
  uri: "mongodb://localhost:27017"
  database: "myapp"
  max_pool_size: 100
  min_pool_size: 10
  connect_timeout: "10s"
  server_selection_timeout: "30s"
```

### 3. Khởi Tạo Ứng Dụng

```go
package main

import (
    "context"
    "log"
    
    "go.fork.vn/app"
    "go.fork.vn/mongodb"
)

func main() {
    // Khởi tạo Fork application
    application := app.New()
    
    // Đăng ký MongoDB provider
    application.RegisterProvider(&mongodb.ServiceProvider{})
    
    // Khởi động application
    if err := application.Boot(); err != nil {
        log.Fatal("Failed to boot application:", err)
    }
    
    // Resolve MongoDB manager
    var mongoManager mongodb.Manager
    if err := application.Container().Resolve(&mongoManager); err != nil {
        log.Fatal("Failed to resolve MongoDB manager:", err)
    }
    
    // Test connection
    ctx := context.Background()
    if err := mongoManager.Ping(ctx); err != nil {
        log.Fatal("MongoDB connection failed:", err)
    }
    
    log.Println("✅ MongoDB connected successfully!")
    
    // Your application logic here
    runApplication(mongoManager)
}
```

## Các Pattern Sử Dụng Phổ Biến

### 1. CRUD Operations

#### Create Documents

```go
func createUser(manager mongodb.Manager, user User) error {
    ctx := context.Background()
    
    // Get collection
    collection, err := manager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return fmt.Errorf("failed to get collection: %w", err)
    }
    
    // Insert document
    result, err := collection.InsertOne(ctx, user)
    if err != nil {
        return fmt.Errorf("failed to insert user: %w", err)
    }
    
    log.Printf("Created user with ID: %v", result.InsertedID)
    return nil
}

// Insert multiple documents
func createUsers(manager mongodb.Manager, users []User) error {
    ctx := context.Background()
    
    collection, err := manager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return err
    }
    
    // Convert to []interface{}
    docs := make([]interface{}, len(users))
    for i, user := range users {
        docs[i] = user
    }
    
    result, err := collection.InsertMany(ctx, docs)
    if err != nil {
        return fmt.Errorf("failed to insert users: %w", err)
    }
    
    log.Printf("Created %d users", len(result.InsertedIDs))
    return nil
}
```

#### Read Documents

```go
import (
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// Find single document
func findUserByID(manager mongodb.Manager, userID string) (*User, error) {
    ctx := context.Background()
    
    collection, err := manager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return nil, err
    }
    
    // Convert string ID to ObjectID
    objectID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return nil, fmt.Errorf("invalid user ID: %w", err)
    }
    
    var user User
    err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, fmt.Errorf("user not found")
        }
        return nil, fmt.Errorf("failed to find user: %w", err)
    }
    
    return &user, nil
}

// Find multiple documents with filters
func findUsersByAge(manager mongodb.Manager, minAge, maxAge int) ([]User, error) {
    ctx := context.Background()
    
    collection, err := manager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return nil, err
    }
    
    // Build filter
    filter := bson.M{
        "age": bson.M{
            "$gte": minAge,
            "$lte": maxAge,
        },
    }
    
    // Set options
    opts := options.Find().SetSort(bson.D{{"name", 1}}).SetLimit(100)
    
    cursor, err := collection.Find(ctx, filter, opts)
    if err != nil {
        return nil, fmt.Errorf("failed to find users: %w", err)
    }
    defer cursor.Close(ctx)
    
    var users []User
    if err = cursor.All(ctx, &users); err != nil {
        return nil, fmt.Errorf("failed to decode users: %w", err)
    }
    
    return users, nil
}

// Pagination example
func findUsersWithPagination(manager mongodb.Manager, page, pageSize int) ([]User, error) {
    ctx := context.Background()
    
    collection, err := manager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return nil, err
    }
    
    skip := (page - 1) * pageSize
    opts := options.Find().
        SetSkip(int64(skip)).
        SetLimit(int64(pageSize)).
        SetSort(bson.D{{"created_at", -1}})
    
    cursor, err := collection.Find(ctx, bson.M{}, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    
    var users []User
    if err = cursor.All(ctx, &users); err != nil {
        return nil, err
    }
    
    return users, nil
}
```

#### Update Documents

```go
// Update single document
func updateUser(manager mongodb.Manager, userID string, updates bson.M) error {
    ctx := context.Background()
    
    collection, err := manager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return err
    }
    
    objectID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return fmt.Errorf("invalid user ID: %w", err)
    }
    
    filter := bson.M{"_id": objectID}
    update := bson.M{
        "$set": updates,
        "$currentDate": bson.M{"updated_at": true},
    }
    
    result, err := collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }
    
    if result.MatchedCount == 0 {
        return fmt.Errorf("user not found")
    }
    
    log.Printf("Updated user %s", userID)
    return nil
}

// Update multiple documents
func updateUsersByStatus(manager mongodb.Manager, oldStatus, newStatus string) error {
    ctx := context.Background()
    
    collection, err := manager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return err
    }
    
    filter := bson.M{"status": oldStatus}
    update := bson.M{
        "$set": bson.M{
            "status": newStatus,
            "updated_at": time.Now(),
        },
    }
    
    result, err := collection.UpdateMany(ctx, filter, update)
    if err != nil {
        return fmt.Errorf("failed to update users: %w", err)
    }
    
    log.Printf("Updated %d users from %s to %s", result.ModifiedCount, oldStatus, newStatus)
    return nil
}

// Upsert example
func upsertUser(manager mongodb.Manager, user User) error {
    ctx := context.Background()
    
    collection, err := manager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return err
    }
    
    filter := bson.M{"email": user.Email}
    update := bson.M{
        "$set": user,
        "$setOnInsert": bson.M{"created_at": time.Now()},
        "$currentDate": bson.M{"updated_at": true},
    }
    
    opts := options.Update().SetUpsert(true)
    
    result, err := collection.UpdateOne(ctx, filter, update, opts)
    if err != nil {
        return fmt.Errorf("failed to upsert user: %w", err)
    }
    
    if result.UpsertedCount > 0 {
        log.Printf("Created new user: %v", result.UpsertedID)
    } else {
        log.Printf("Updated existing user")
    }
    
    return nil
}
```

#### Delete Documents

```go
// Delete single document
func deleteUser(manager mongodb.Manager, userID string) error {
    ctx := context.Background()
    
    collection, err := manager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return err
    }
    
    objectID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return fmt.Errorf("invalid user ID: %w", err)
    }
    
    result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
    if err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }
    
    if result.DeletedCount == 0 {
        return fmt.Errorf("user not found")
    }
    
    log.Printf("Deleted user %s", userID)
    return nil
}

// Delete multiple documents
func deleteInactiveUsers(manager mongodb.Manager, days int) error {
    ctx := context.Background()
    
    collection, err := manager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return err
    }
    
    cutoff := time.Now().AddDate(0, 0, -days)
    filter := bson.M{
        "last_login": bson.M{"$lt": cutoff},
        "status": "inactive",
    }
    
    result, err := collection.DeleteMany(ctx, filter)
    if err != nil {
        return fmt.Errorf("failed to delete inactive users: %w", err)
    }
    
    log.Printf("Deleted %d inactive users", result.DeletedCount)
    return nil
}
```

### 2. Aggregation Pipeline

```go
import "go.mongodb.org/mongo-driver/mongo"

// Group users by age and count
func getUserCountByAge(manager mongodb.Manager) ([]bson.M, error) {
    ctx := context.Background()
    
    collection, err := manager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return nil, err
    }
    
    pipeline := mongo.Pipeline{
        {{"$group", bson.D{
            {"_id", "$age"},
            {"count", bson.D{{"$sum", 1}}},
            {"avg_score", bson.D{{"$avg", "$score"}}},
        }}},
        {{"$sort", bson.D{{"_id", 1}}}},
    }
    
    cursor, err := collection.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, fmt.Errorf("aggregation failed: %w", err)
    }
    defer cursor.Close(ctx)
    
    var results []bson.M
    if err = cursor.All(ctx, &results); err != nil {
        return nil, fmt.Errorf("failed to decode results: %w", err)
    }
    
    return results, nil
}

// Complex aggregation with lookup
func getUsersWithOrders(manager mongodb.Manager) ([]bson.M, error) {
    ctx := context.Background()
    
    collection, err := manager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return nil, err
    }
    
    pipeline := mongo.Pipeline{
        {{"$lookup", bson.D{
            {"from", "orders"},
            {"localField", "_id"},
            {"foreignField", "user_id"},
            {"as", "orders"},
        }}},
        {{"$addFields", bson.D{
            {"total_orders", bson.D{{"$size", "$orders"}}},
            {"total_amount", bson.D{{"$sum", "$orders.amount"}}},
        }}},
        {{"$match", bson.D{
            {"total_orders", bson.D{{"$gt", 0}}},
        }}},
        {{"$sort", bson.D{{"total_amount", -1}}}},
        {{"$limit", 100}},
    }
    
    cursor, err := collection.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, fmt.Errorf("aggregation failed: %w", err)
    }
    defer cursor.Close(ctx)
    
    var results []bson.M
    if err = cursor.All(ctx, &results); err != nil {
        return nil, err
    }
    
    return results, nil
}
```

### 3. Transactions

```go
import "go.mongodb.org/mongo-driver/mongo"

func transferMoney(manager mongodb.Manager, fromUserID, toUserID string, amount float64) error {
    ctx := context.Background()
    
    client, err := manager.GetClient(ctx)
    if err != nil {
        return err
    }
    
    // Start session
    session, err := client.StartSession()
    if err != nil {
        return fmt.Errorf("failed to start session: %w", err)
    }
    defer session.EndSession(ctx)
    
    // Define transaction function
    callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
        usersCollection, err := manager.GetCollection(sessCtx, "myapp", "users")
        if err != nil {
            return nil, err
        }
        
        // Deduct from sender
        fromFilter := bson.M{"_id": fromUserID, "balance": bson.M{"$gte": amount}}
        fromUpdate := bson.M{"$inc": bson.M{"balance": -amount}}
        
        fromResult, err := usersCollection.UpdateOne(sessCtx, fromFilter, fromUpdate)
        if err != nil {
            return nil, fmt.Errorf("failed to deduct from sender: %w", err)
        }
        
        if fromResult.ModifiedCount == 0 {
            return nil, fmt.Errorf("insufficient balance or sender not found")
        }
        
        // Add to receiver
        toFilter := bson.M{"_id": toUserID}
        toUpdate := bson.M{"$inc": bson.M{"balance": amount}}
        
        toResult, err := usersCollection.UpdateOne(sessCtx, toFilter, toUpdate)
        if err != nil {
            return nil, fmt.Errorf("failed to add to receiver: %w", err)
        }
        
        if toResult.ModifiedCount == 0 {
            return nil, fmt.Errorf("receiver not found")
        }
        
        // Record transaction
        transactionsCollection, err := manager.GetCollection(sessCtx, "myapp", "transactions")
        if err != nil {
            return nil, err
        }
        
        transaction := bson.M{
            "from_user_id": fromUserID,
            "to_user_id":   toUserID,
            "amount":       amount,
            "timestamp":    time.Now(),
            "status":       "completed",
        }
        
        _, err = transactionsCollection.InsertOne(sessCtx, transaction)
        if err != nil {
            return nil, fmt.Errorf("failed to record transaction: %w", err)
        }
        
        return nil, nil
    }
    
    // Execute transaction with retry
    _, err = session.WithTransaction(ctx, callback)
    if err != nil {
        return fmt.Errorf("transaction failed: %w", err)
    }
    
    log.Printf("Successfully transferred %.2f from %s to %s", amount, fromUserID, toUserID)
    return nil
}
```

### 4. Indexes

```go
import (
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func createIndexes(manager mongodb.Manager) error {
    ctx := context.Background()
    
    collection, err := manager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return err
    }
    
    // Single field index
    emailIndex := mongo.IndexModel{
        Keys:    bson.D{{"email", 1}},
        Options: options.Index().SetUnique(true).SetBackground(true),
    }
    
    // Compound index
    nameAgeIndex := mongo.IndexModel{
        Keys: bson.D{
            {"name", 1},
            {"age", -1},
        },
        Options: options.Index().SetBackground(true),
    }
    
    // Text index
    textIndex := mongo.IndexModel{
        Keys: bson.D{
            {"name", "text"},
            {"description", "text"},
        },
        Options: options.Index().SetBackground(true),
    }
    
    // TTL index
    ttlIndex := mongo.IndexModel{
        Keys:    bson.D{{"expires_at", 1}},
        Options: options.Index().SetExpireAfterSeconds(0).SetBackground(true),
    }
    
    indexes := []mongo.IndexModel{emailIndex, nameAgeIndex, textIndex, ttlIndex}
    
    names, err := collection.Indexes().CreateMany(ctx, indexes)
    if err != nil {
        return fmt.Errorf("failed to create indexes: %w", err)
    }
    
    log.Printf("Created indexes: %v", names)
    return nil
}

// List existing indexes
func listIndexes(manager mongodb.Manager) error {
    ctx := context.Background()
    
    collection, err := manager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return err
    }
    
    cursor, err := collection.Indexes().List(ctx)
    if err != nil {
        return fmt.Errorf("failed to list indexes: %w", err)
    }
    defer cursor.Close(ctx)
    
    var indexes []bson.M
    if err = cursor.All(ctx, &indexes); err != nil {
        return fmt.Errorf("failed to decode indexes: %w", err)
    }
    
    for _, index := range indexes {
        log.Printf("Index: %+v", index)
    }
    
    return nil
}
```

## Cấu Hình Nâng Cao

### 1. SSL/TLS Configuration

```yaml
# configs/app.yaml
mongodb:
  uri: "mongodb://localhost:27017"
  database: "myapp"
  ssl:
    enabled: true
    ca_file: "/etc/ssl/certs/mongodb-ca.pem"
    certificate_file: "/etc/ssl/certs/mongodb-client.pem"
    private_key_file: "/etc/ssl/private/mongodb-client-key.pem"
    insecure_skip_verify: false
```

### 2. Replica Set Configuration

```yaml
mongodb:
  uri: "mongodb://host1:27017,host2:27017,host3:27017/?replicaSet=rs0"
  database: "myapp"
  max_pool_size: 200
  min_pool_size: 20
  server_selection_timeout: "60s"
```

### 3. Authentication

```yaml
mongodb:
  uri: "mongodb://localhost:27017"
  database: "myapp"
  auth:
    username: "app_user"
    password: "secure_password"
    auth_db: "admin"
```

### 4. Environment Variables

```bash
# Connection
export MONGODB_URI="mongodb://localhost:27017"
export MONGODB_DATABASE="myapp"

# Pool settings
export MONGODB_MAX_POOL_SIZE=100
export MONGODB_MIN_POOL_SIZE=10
export MONGODB_MAX_CONN_IDLE_TIME="30s"

# Timeouts
export MONGODB_SERVER_SELECTION_TIMEOUT="30s"
export MONGODB_CONNECT_TIMEOUT="10s"
export MONGODB_SOCKET_TIMEOUT="30s"

# SSL
export MONGODB_SSL_ENABLED=true
export MONGODB_SSL_CA_FILE="/path/to/ca.pem"
export MONGODB_SSL_CERTIFICATE_FILE="/path/to/cert.pem"
export MONGODB_SSL_PRIVATE_KEY_FILE="/path/to/key.pem"

# Authentication
export MONGODB_AUTH_USERNAME="admin"
export MONGODB_AUTH_PASSWORD="password123"
export MONGODB_AUTH_DB="admin"
```

## Testing Patterns

### 1. Unit Tests với Mocks

```go
package service_test

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "go.mongodb.org/mongo-driver/mongo"
    
    "go.fork.vn/mongodb/mocks"
    "your-app/service"
)

func TestUserService_CreateUser(t *testing.T) {
    // Setup mock
    mockManager := &mocks.MockManager{}
    mockCollection := &mongo.Collection{} // mock này cần setup riêng
    
    // Setup expectations
    mockManager.On("GetCollection", mock.Anything, "myapp", "users").
        Return(mockCollection, nil)
    
    // Test service
    userService := service.NewUserService(mockManager)
    
    user := &User{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    err := userService.CreateUser(context.Background(), user)
    
    // Assertions
    assert.NoError(t, err)
    mockManager.AssertExpectations(t)
}
```

### 2. Integration Tests

```go
func TestMongoDBIntegration(t *testing.T) {
    // Skip if not integration test
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    // Setup test database
    config := &mongodb.Config{
        URI:      "mongodb://localhost:27017",
        Database: "test_db",
    }
    
    manager := mongodb.NewManager(config)
    ctx := context.Background()
    
    // Test connection
    err := manager.Ping(ctx)
    assert.NoError(t, err)
    
    // Test operations
    collection, err := manager.GetCollection(ctx, "test_db", "test_collection")
    assert.NoError(t, err)
    
    // Cleanup
    defer func() {
        collection.Drop(ctx)
        manager.Disconnect(ctx)
    }()
    
    // Test insert
    doc := bson.M{"name": "test", "value": 123}
    result, err := collection.InsertOne(ctx, doc)
    assert.NoError(t, err)
    assert.NotNil(t, result.InsertedID)
}
```

### 3. Test Helpers

```go
package testutil

import (
    "context"
    "fmt"
    "testing"
    
    "go.fork.vn/mongodb"
)

func SetupTestDB(t *testing.T) mongodb.Manager {
    config := &mongodb.Config{
        URI:      getTestMongoURI(),
        Database: fmt.Sprintf("test_%s", t.Name()),
    }
    
    manager := mongodb.NewManager(config)
    
    // Test connection
    ctx := context.Background()
    if err := manager.Ping(ctx); err != nil {
        t.Fatalf("Failed to connect to test MongoDB: %v", err)
    }
    
    return manager
}

func CleanupTestDB(t *testing.T, manager mongodb.Manager) {
    ctx := context.Background()
    
    // Drop test database
    client, err := manager.GetClient(ctx)
    if err != nil {
        t.Logf("Failed to get client for cleanup: %v", err)
        return
    }
    
    dbName := fmt.Sprintf("test_%s", t.Name())
    if err := client.Database(dbName).Drop(ctx); err != nil {
        t.Logf("Failed to drop test database: %v", err)
    }
    
    // Disconnect
    if err := manager.Disconnect(ctx); err != nil {
        t.Logf("Failed to disconnect: %v", err)
    }
}

func getTestMongoURI() string {
    uri := os.Getenv("TEST_MONGODB_URI")
    if uri == "" {
        uri = "mongodb://localhost:27017"
    }
    return uri
}
```

## Performance Optimization

### 1. Connection Pool Tuning

```yaml
mongodb:
  max_pool_size: 100        # Điều chỉnh theo concurrent requests
  min_pool_size: 10         # Maintain minimum connections
  max_conn_idle_time: "30s" # Release idle connections
  server_selection_timeout: "30s"
  connect_timeout: "10s"
  socket_timeout: "30s"
```

### 2. Query Optimization

```go
// Use projection để giảm data transfer
func findUsersProjection(manager mongodb.Manager) ([]bson.M, error) {
    ctx := context.Background()
    
    collection, err := manager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return nil, err
    }
    
    // Only select needed fields
    opts := options.Find().SetProjection(bson.M{
        "name":  1,
        "email": 1,
        "_id":   0,
    })
    
    cursor, err := collection.Find(ctx, bson.M{}, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)
    
    var results []bson.M
    if err = cursor.All(ctx, &results); err != nil {
        return nil, err
    }
    
    return results, nil
}

// Use hints cho query optimization
func findWithHint(manager mongodb.Manager) error {
    ctx := context.Background()
    
    collection, err := manager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return err
    }
    
    opts := options.Find().SetHint(bson.D{{"email", 1}})
    
    cursor, err := collection.Find(ctx, bson.M{"email": "john@example.com"}, opts)
    if err != nil {
        return err
    }
    defer cursor.Close(ctx)
    
    return nil
}
```

### 3. Bulk Operations

```go
// Bulk insert
func bulkInsertUsers(manager mongodb.Manager, users []User) error {
    ctx := context.Background()
    
    collection, err := manager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return err
    }
    
    // Convert to []interface{}
    docs := make([]interface{}, len(users))
    for i, user := range users {
        docs[i] = user
    }
    
    // Use bulk insert
    result, err := collection.InsertMany(ctx, docs)
    if err != nil {
        return fmt.Errorf("bulk insert failed: %w", err)
    }
    
    log.Printf("Inserted %d documents", len(result.InsertedIDs))
    return nil
}

// Bulk write operations
func bulkWriteOperations(manager mongodb.Manager) error {
    ctx := context.Background()
    
    collection, err := manager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return err
    }
    
    operations := []mongo.WriteModel{
        mongo.NewInsertOneModel().SetDocument(bson.M{"name": "Alice", "age": 25}),
        mongo.NewUpdateOneModel().SetFilter(bson.M{"name": "Bob"}).
            SetUpdate(bson.M{"$set": bson.M{"age": 30}}),
        mongo.NewDeleteOneModel().SetFilter(bson.M{"name": "Charlie"}),
    }
    
    opts := options.BulkWrite().SetOrdered(false)
    result, err := collection.BulkWrite(ctx, operations, opts)
    if err != nil {
        return fmt.Errorf("bulk write failed: %w", err)
    }
    
    log.Printf("Bulk write results: %+v", result)
    return nil
}
```

## Error Handling Best Practices

### 1. Context Handling

```go
func operationWithTimeout(manager mongodb.Manager) error {
    // Create context với timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    collection, err := manager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return fmt.Errorf("failed to get collection: %w", err)
    }
    
    // Operation sẽ timeout sau 30 giây
    _, err = collection.InsertOne(ctx, bson.M{"name": "test"})
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            return fmt.Errorf("operation timed out")
        }
        return fmt.Errorf("insert failed: %w", err)
    }
    
    return nil
}
```

### 2. Retry Logic

```go
import "time"

func operationWithRetry(manager mongodb.Manager, maxRetries int) error {
    ctx := context.Background()
    
    for attempt := 1; attempt <= maxRetries; attempt++ {
        err := performOperation(manager, ctx)
        if err == nil {
            return nil
        }
        
        // Check if error is retryable
        if !isRetryableError(err) {
            return fmt.Errorf("non-retryable error: %w", err)
        }
        
        if attempt < maxRetries {
            // Exponential backoff
            backoff := time.Duration(attempt) * time.Second
            log.Printf("Attempt %d failed, retrying in %v: %v", attempt, backoff, err)
            time.Sleep(backoff)
        }
    }
    
    return fmt.Errorf("operation failed after %d attempts", maxRetries)
}

func isRetryableError(err error) bool {
    if err == nil {
        return false
    }
    
    // Check for specific MongoDB errors that are retryable
    if mongo.IsTimeout(err) || mongo.IsNetworkError(err) {
        return true
    }
    
    return false
}

func performOperation(manager mongodb.Manager, ctx context.Context) error {
    collection, err := manager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        return err
    }
    
    _, err = collection.InsertOne(ctx, bson.M{"name": "test"})
    return err
}
```

### 3. Graceful Shutdown

```go
import (
    "os"
    "os/signal"
    "syscall"
)

func runApplicationWithGracefulShutdown(manager mongodb.Manager) {
    // Setup signal handling
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    // Run application
    go func() {
        log.Println("Application starting...")
        // Your application logic here
        runBusinessLogic(manager)
    }()
    
    // Wait for signal
    sig := <-sigChan
    log.Printf("Received signal: %v", sig)
    
    // Graceful shutdown
    log.Println("Shutting down gracefully...")
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := manager.Disconnect(ctx); err != nil {
        log.Printf("Error disconnecting from MongoDB: %v", err)
    } else {
        log.Println("✅ MongoDB disconnected successfully")
    }
    
    log.Println("Application stopped")
}
```

## Monitoring và Logging

### 1. Connection Monitoring

```go
func monitorConnections(manager mongodb.Manager) {
    ctx := context.Background()
    
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if manager.IsConnected(ctx) {
                log.Println("✅ MongoDB connection healthy")
            } else {
                log.Println("❌ MongoDB connection lost")
                
                // Attempt reconnection
                if err := manager.Ping(ctx); err != nil {
                    log.Printf("Reconnection failed: %v", err)
                } else {
                    log.Println("✅ MongoDB reconnected")
                }
            }
        }
    }
}
```

### 2. Performance Metrics

```go
import "time"

func operationWithMetrics(manager mongodb.Manager) error {
    start := time.Now()
    ctx := context.Background()
    
    defer func() {
        duration := time.Since(start)
        log.Printf("Operation completed in %v", duration)
        
        // Send metrics to monitoring system
        // metrics.Timer("mongodb.operation.duration", duration)
    }()
    
    collection, err := manager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        // metrics.Counter("mongodb.errors", 1, []string{"type:connection"})
        return err
    }
    
    _, err = collection.InsertOne(ctx, bson.M{"name": "test"})
    if err != nil {
        // metrics.Counter("mongodb.errors", 1, []string{"type:operation"})
        return err
    }
    
    // metrics.Counter("mongodb.operations", 1, []string{"type:insert"})
    return nil
}
```

---

> 📘 **Tip**: Để biết thêm chi tiết về các tính năng nâng cao, tham khảo [API Reference](reference.md) và [Overview](overview.md).
