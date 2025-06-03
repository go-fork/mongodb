# API Reference - MongoDB Provider v0.1.1

## Package Overview

Package `go.fork.vn/mongodb` cung cấp MongoDB integration cho Fork Framework với high-performance connection management và type-safe operations.

```go
import "go.fork.vn/mongodb"
```

## Core Interfaces

### Manager Interface

Manager là interface chính để tương tác với MongoDB instances.

```go
type Manager interface {
    GetClient(ctx context.Context) (*mongo.Client, error)
    GetDatabase(ctx context.Context, name string) (*mongo.Database, error)
    GetCollection(ctx context.Context, database, collection string) (*mongo.Collection, error)
    Ping(ctx context.Context) error
    Disconnect(ctx context.Context) error
    IsConnected(ctx context.Context) bool
    GetConnectionString() string
}
```

#### Methods

##### GetClient

```go
GetClient(ctx context.Context) (*mongo.Client, error)
```

Trả về MongoDB client instance từ connection pool.

**Parameters:**
- `ctx context.Context`: Context cho operation (support cancellation/timeout)

**Returns:**
- `*mongo.Client`: MongoDB client instance
- `error`: Lỗi nếu không thể lấy client

**Example:**
```go
client, err := manager.GetClient(ctx)
if err != nil {
    return fmt.Errorf("failed to get client: %w", err)
}
defer client.Disconnect(ctx)
```

**Errors:**
- `ConnectionError`: Khi không thể kết nối đến MongoDB
- `ContextError`: Khi context bị cancelled hoặc timeout

##### GetDatabase

```go
GetDatabase(ctx context.Context, name string) (*mongo.Database, error)
```

Trả về database instance với tên được chỉ định.

**Parameters:**
- `ctx context.Context`: Context cho operation
- `name string`: Tên database

**Returns:**
- `*mongo.Database`: Database instance
- `error`: Lỗi nếu không thể lấy database

**Example:**
```go
db, err := manager.GetDatabase(ctx, "myapp")
if err != nil {
    return fmt.Errorf("failed to get database: %w", err)
}
```

**Behavior:**
- Nếu `name` trống, sử dụng default database từ configuration
- Database sẽ được tạo tự động nếu chưa tồn tại (khi có operations)

##### GetCollection

```go
GetCollection(ctx context.Context, database, collection string) (*mongo.Collection, error)
```

Trả về collection instance từ database được chỉ định.

**Parameters:**
- `ctx context.Context`: Context cho operation
- `database string`: Tên database
- `collection string`: Tên collection

**Returns:**
- `*mongo.Collection`: Collection instance
- `error`: Lỗi nếu không thể lấy collection

**Example:**
```go
collection, err := manager.GetCollection(ctx, "myapp", "users")
if err != nil {
    return fmt.Errorf("failed to get collection: %w", err)
}
```

##### Ping

```go
Ping(ctx context.Context) error
```

Kiểm tra kết nối đến MongoDB bằng cách gửi ping command.

**Parameters:**
- `ctx context.Context`: Context cho operation

**Returns:**
- `error`: Lỗi nếu ping thất bại

**Example:**
```go
if err := manager.Ping(ctx); err != nil {
    log.Printf("MongoDB ping failed: %v", err)
    return err
}
```

**Timeout:** Sử dụng timeout từ context hoặc server selection timeout từ config.

##### IsConnected

```go
IsConnected(ctx context.Context) bool
```

Kiểm tra trạng thái kết nối hiện tại.

**Parameters:**
- `ctx context.Context`: Context cho operation

**Returns:**
- `bool`: `true` nếu đã kết nối, `false` nếu ngược lại

**Example:**
```go
if !manager.IsConnected(ctx) {
    return errors.New("MongoDB not connected")
}
```

**Note:** Method này perform một ping ngắn để verify connection health.

##### Disconnect

```go
Disconnect(ctx context.Context) error
```

Đóng tất cả connections đến MongoDB.

**Parameters:**
- `ctx context.Context`: Context cho operation

**Returns:**
- `error`: Lỗi nếu disconnect thất bại

**Example:**
```go
defer func() {
    if err := manager.Disconnect(ctx); err != nil {
        log.Printf("Error disconnecting: %v", err)
    }
}()
```

**Behavior:**
- Đóng tất cả active connections
- Chờ ongoing operations complete (respect context timeout)
- Cleanup connection pool resources

##### GetConnectionString

```go
GetConnectionString() string
```

Trả về connection string đang được sử dụng (đã sanitized, không chứa credentials).

**Returns:**
- `string`: Connection string (không có password)

**Example:**
```go
connStr := manager.GetConnectionString()
log.Printf("Connected to: %s", connStr)
```

## Configuration Types

### Config Struct

```go
type Config struct {
    URI                    string        `mapstructure:"uri" yaml:"uri"`
    Database               string        `mapstructure:"database" yaml:"database"`
    MaxPoolSize            *uint64       `mapstructure:"max_pool_size" yaml:"max_pool_size"`
    MinPoolSize            *uint64       `mapstructure:"min_pool_size" yaml:"min_pool_size"`
    MaxConnIdleTime        time.Duration `mapstructure:"max_conn_idle_time" yaml:"max_conn_idle_time"`
    ServerSelectionTimeout time.Duration `mapstructure:"server_selection_timeout" yaml:"server_selection_timeout"`
    ConnectTimeout         time.Duration `mapstructure:"connect_timeout" yaml:"connect_timeout"`
    SocketTimeout          time.Duration `mapstructure:"socket_timeout" yaml:"socket_timeout"`
    SSL                    SSLConfig     `mapstructure:"ssl" yaml:"ssl"`
    Auth                   AuthConfig    `mapstructure:"auth" yaml:"auth"`
}
```

#### Fields

##### URI
```go
URI string `mapstructure:"uri" yaml:"uri"`
```
MongoDB connection URI.

**Format:** `mongodb://[username:password@]host1[:port1][,...hostN[:portN]][/[defaultauthdb][?options]]`

**Example:**
```yaml
uri: "mongodb://localhost:27017"
uri: "mongodb://user:pass@localhost:27017/mydb"
uri: "mongodb://host1:27017,host2:27017,host3:27017"
```

##### Database
```go
Database string `mapstructure:"database" yaml:"database"`
```
Default database name.

**Example:**
```yaml
database: "myapp"
```

##### MaxPoolSize
```go
MaxPoolSize *uint64 `mapstructure:"max_pool_size" yaml:"max_pool_size"`
```
Số connection tối đa trong pool.

**Default:** `100`
**Range:** `1-1000`

**Example:**
```yaml
max_pool_size: 100
```

##### MinPoolSize
```go
MinPoolSize *uint64 `mapstructure:"min_pool_size" yaml:"min_pool_size"`
```
Số connection tối thiểu maintain trong pool.

**Default:** `10`
**Range:** `0-MaxPoolSize`

**Example:**
```yaml
min_pool_size: 10
```

##### MaxConnIdleTime
```go
MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time" yaml:"max_conn_idle_time"`
```
Thời gian tối đa connection có thể idle.

**Default:** `30s`
**Format:** Go duration string

**Example:**
```yaml
max_conn_idle_time: "30s"
max_conn_idle_time: "5m"
```

##### ServerSelectionTimeout
```go
ServerSelectionTimeout time.Duration `mapstructure:"server_selection_timeout" yaml:"server_selection_timeout"`
```
Timeout cho server selection.

**Default:** `30s`

**Example:**
```yaml
server_selection_timeout: "30s"
```

##### ConnectTimeout
```go
ConnectTimeout time.Duration `mapstructure:"connect_timeout" yaml:"connect_timeout"`
```
Timeout cho connection establishment.

**Default:** `10s`

**Example:**
```yaml
connect_timeout: "10s"
```

##### SocketTimeout
```go
SocketTimeout time.Duration `mapstructure:"socket_timeout" yaml:"socket_timeout"`
```
Timeout cho socket operations.

**Default:** `30s`

**Example:**
```yaml
socket_timeout: "30s"
```

### SSLConfig Struct

```go
type SSLConfig struct {
    Enabled            bool   `mapstructure:"enabled" yaml:"enabled"`
    CAFile             string `mapstructure:"ca_file" yaml:"ca_file"`
    CertificateFile    string `mapstructure:"certificate_file" yaml:"certificate_file"`
    PrivateKeyFile     string `mapstructure:"private_key_file" yaml:"private_key_file"`
    InsecureSkipVerify bool   `mapstructure:"insecure_skip_verify" yaml:"insecure_skip_verify"`
}
```

#### Fields

##### Enabled
```go
Enabled bool `mapstructure:"enabled" yaml:"enabled"`
```
Bật/tắt SSL/TLS encryption.

**Default:** `false`

##### CAFile
```go
CAFile string `mapstructure:"ca_file" yaml:"ca_file"`
```
Đường dẫn đến CA certificate file.

**Example:**
```yaml
ca_file: "/path/to/ca.pem"
```

##### CertificateFile
```go
CertificateFile string `mapstructure:"certificate_file" yaml:"certificate_file"`
```
Đường dẫn đến client certificate file.

**Example:**
```yaml
certificate_file: "/path/to/client.pem"
```

##### PrivateKeyFile
```go
PrivateKeyFile string `mapstructure:"private_key_file" yaml:"private_key_file"`
```
Đường dẫn đến private key file.

**Example:**
```yaml
private_key_file: "/path/to/client-key.pem"
```

##### InsecureSkipVerify
```go
InsecureSkipVerify bool `mapstructure:"insecure_skip_verify" yaml:"insecure_skip_verify"`
```
Bỏ qua certificate verification (chỉ dùng cho development).

**Default:** `false`
**Security Warning:** Không sử dụng trong production.

### AuthConfig Struct

```go
type AuthConfig struct {
    Username string `mapstructure:"username" yaml:"username"`
    Password string `mapstructure:"password" yaml:"password"`
    AuthDB   string `mapstructure:"auth_db" yaml:"auth_db"`
}
```

#### Fields

##### Username
```go
Username string `mapstructure:"username" yaml:"username"`
```
Username cho authentication.

##### Password
```go
Password string `mapstructure:"password" yaml:"password"`
```
Password cho authentication.

**Security Note:** Nên sử dụng environment variables cho password.

##### AuthDB
```go
AuthDB string `mapstructure:"auth_db" yaml:"auth_db"`
```
Database sử dụng cho authentication.

**Default:** `"admin"`

## Service Provider

### ServiceProvider Struct

```go
type ServiceProvider struct {
    configManager config.Manager
}
```

Implement `di.ServiceProvider` interface để tích hợp với DI container.

#### Methods

##### NewServiceProvider

```go
func NewServiceProvider() *ServiceProvider
```

Tạo ServiceProvider instance mới.

**Returns:**
- `*ServiceProvider`: ServiceProvider instance

**Example:**
```go
provider := mongodb.NewServiceProvider()
```

##### Register

```go
func (p *ServiceProvider) Register(container di.Container) error
```

Đăng ký MongoDB services vào DI container.

**Parameters:**
- `container di.Container`: DI container instance

**Returns:**
- `error`: Lỗi nếu registration thất bại

**Behavior:**
- Đăng ký `Manager` interface as singleton
- Tự động resolve config dependencies

##### Boot

```go
func (p *ServiceProvider) Boot(container di.Container) error
```

Bootstrap MongoDB services sau khi container được build.

**Parameters:**
- `container di.Container`: DI container instance

**Returns:**
- `error`: Lỗi nếu boot thất bại

## Factory Functions

### NewManager

```go
func NewManager(config *Config) Manager
```

Tạo Manager instance mới với configuration.

**Parameters:**
- `config *Config`: MongoDB configuration

**Returns:**
- `Manager`: Manager instance

**Example:**
```go
config := &mongodb.Config{
    URI:      "mongodb://localhost:27017",
    Database: "myapp",
}
manager := mongodb.NewManager(config)
```

## Environment Variables

Package hỗ trợ cấu hình thông qua environment variables:

| Variable | Type | Description | Default |
|----------|------|-------------|---------|
| `MONGODB_URI` | string | Connection URI | `mongodb://localhost:27017` |
| `MONGODB_DATABASE` | string | Default database | `""` |
| `MONGODB_MAX_POOL_SIZE` | uint64 | Max pool size | `100` |
| `MONGODB_MIN_POOL_SIZE` | uint64 | Min pool size | `10` |
| `MONGODB_MAX_CONN_IDLE_TIME` | duration | Max idle time | `30s` |
| `MONGODB_SERVER_SELECTION_TIMEOUT` | duration | Server selection timeout | `30s` |
| `MONGODB_CONNECT_TIMEOUT` | duration | Connect timeout | `10s` |
| `MONGODB_SOCKET_TIMEOUT` | duration | Socket timeout | `30s` |
| `MONGODB_SSL_ENABLED` | bool | Enable SSL | `false` |
| `MONGODB_SSL_CA_FILE` | string | CA file path | `""` |
| `MONGODB_SSL_CERTIFICATE_FILE` | string | Certificate file path | `""` |
| `MONGODB_SSL_PRIVATE_KEY_FILE` | string | Private key file path | `""` |
| `MONGODB_SSL_INSECURE_SKIP_VERIFY` | bool | Skip SSL verification | `false` |
| `MONGODB_AUTH_USERNAME` | string | Auth username | `""` |
| `MONGODB_AUTH_PASSWORD` | string | Auth password | `""` |
| `MONGODB_AUTH_DB` | string | Auth database | `admin` |

## Error Types

### ConnectionError

```go
type ConnectionError struct {
    Err error
}
```

Lỗi liên quan đến MongoDB connection.

#### Methods

```go
func (e *ConnectionError) Error() string
func (e *ConnectionError) Unwrap() error
func (e *ConnectionError) Temporary() bool // returns true
```

### AuthenticationError

```go
type AuthenticationError struct {
    Err error
}
```

Lỗi authentication với MongoDB.

#### Methods

```go
func (e *AuthenticationError) Error() string
func (e *AuthenticationError) Unwrap() error
func (e *AuthenticationError) Temporary() bool // returns false
```

### ConfigurationError

```go
type ConfigurationError struct {
    Field string
    Value interface{}
    Err   error
}
```

Lỗi configuration validation.

#### Methods

```go
func (e *ConfigurationError) Error() string
func (e *ConfigurationError) Unwrap() error
```

## Constants

```go
const (
    // Default configuration values
    DefaultMaxPoolSize            = 100
    DefaultMinPoolSize            = 10
    DefaultMaxConnIdleTime        = 30 * time.Second
    DefaultServerSelectionTimeout = 30 * time.Second
    DefaultConnectTimeout         = 10 * time.Second
    DefaultSocketTimeout          = 30 * time.Second
    DefaultAuthDB                 = "admin"
    
    // Environment variable names
    EnvURI                    = "MONGODB_URI"
    EnvDatabase               = "MONGODB_DATABASE"
    EnvMaxPoolSize            = "MONGODB_MAX_POOL_SIZE"
    EnvMinPoolSize            = "MONGODB_MIN_POOL_SIZE"
    EnvMaxConnIdleTime        = "MONGODB_MAX_CONN_IDLE_TIME"
    EnvServerSelectionTimeout = "MONGODB_SERVER_SELECTION_TIMEOUT"
    EnvConnectTimeout         = "MONGODB_CONNECT_TIMEOUT"
    EnvSocketTimeout          = "MONGODB_SOCKET_TIMEOUT"
    
    // SSL Environment variables
    EnvSSLEnabled            = "MONGODB_SSL_ENABLED"
    EnvSSLCAFile             = "MONGODB_SSL_CA_FILE"
    EnvSSLCertificateFile    = "MONGODB_SSL_CERTIFICATE_FILE"
    EnvSSLPrivateKeyFile     = "MONGODB_SSL_PRIVATE_KEY_FILE"
    EnvSSLInsecureSkipVerify = "MONGODB_SSL_INSECURE_SKIP_VERIFY"
    
    // Auth Environment variables
    EnvAuthUsername = "MONGODB_AUTH_USERNAME"
    EnvAuthPassword = "MONGODB_AUTH_PASSWORD"
    EnvAuthDB       = "MONGODB_AUTH_DB"
)
```

## Example Usage

### Basic Usage

```go
package main

import (
    "context"
    "log"
    
    "go.fork.vn/app"
    "go.fork.vn/mongodb"
)

func main() {
    // Setup application
    application := app.New()
    application.RegisterProvider(&mongodb.ServiceProvider{})
    
    if err := application.Boot(); err != nil {
        log.Fatal(err)
    }
    
    // Resolve manager
    var manager mongodb.Manager
    if err := application.Container().Resolve(&manager); err != nil {
        log.Fatal(err)
    }
    
    // Use manager
    ctx := context.Background()
    collection, err := manager.GetCollection(ctx, "myapp", "users")
    if err != nil {
        log.Fatal(err)
    }
    
    // Perform operations
    result, err := collection.InsertOne(ctx, map[string]interface{}{
        "name": "John Doe",
        "email": "john@example.com",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Inserted document with ID: %v", result.InsertedID)
}
```

### Advanced Configuration

```go
config := &mongodb.Config{
    URI:      "mongodb://user:pass@host1:27017,host2:27017/mydb?replicaSet=rs0",
    Database: "myapp",
    MaxPoolSize: &[]uint64{200}[0],
    MinPoolSize: &[]uint64{20}[0],
    MaxConnIdleTime: 5 * time.Minute,
    ServerSelectionTimeout: 60 * time.Second,
    ConnectTimeout: 30 * time.Second,
    SocketTimeout: 60 * time.Second,
    SSL: mongodb.SSLConfig{
        Enabled:            true,
        CAFile:             "/etc/ssl/certs/mongodb-ca.pem",
        CertificateFile:    "/etc/ssl/certs/mongodb-client.pem",
        PrivateKeyFile:     "/etc/ssl/private/mongodb-client-key.pem",
        InsecureSkipVerify: false,
    },
    Auth: mongodb.AuthConfig{
        Username: "app_user",
        Password: "secure_password",
        AuthDB:   "admin",
    },
}

manager := mongodb.NewManager(config)
```

---

> 📘 **API Stability**: API này được thiết kế để backward compatible. Breaking changes sẽ được thông báo trước ít nhất 1 major version.
