package mongodb_test

import (
	"reflect"
	"testing"
	"time"

	"go.fork.vn/mongodb"
)

func TestDefaultConfig(t *testing.T) {
	expected := &mongodb.Config{
		URI:                    "mongodb://localhost:27017",
		Database:               "myapp",
		AppName:                "golang-app",
		MaxPoolSize:            100,
		MinPoolSize:            5,
		MaxConnecting:          10,
		MaxConnIdleTime:        600000,
		ConnectTimeout:         30000,
		ServerSelectionTimeout: 30000,
		SocketTimeout:          0,
		HeartbeatInterval:      10000,
		LocalThreshold:         15000,
		Timeout:                30000,
		TLS: mongodb.TLSConfig{
			Enabled:            false,
			InsecureSkipVerify: false,
		},
		Auth: mongodb.AuthConfig{
			AuthSource:              "admin",
			AuthMechanism:           "SCRAM-SHA-256",
			AuthMechanismProperties: make(map[string]string),
		},
		ReadPreference: mongodb.ReadPreferenceConfig{
			Mode:         "primary",
			TagSets:      []map[string]string{},
			MaxStaleness: 90,
			HedgeEnabled: false,
		},
		ReadConcern: mongodb.ReadConcernConfig{
			Level: "majority",
		},
		WriteConcern: mongodb.WriteConcernConfig{
			W:        "majority",
			Journal:  true,
			WTimeout: 30000,
		},
		RetryWrites:  true,
		RetryReads:   true,
		Compressors:  []string{},
		ZlibLevel:    6,
		ZstdLevel:    6,
		Direct:       false,
		LoadBalanced: false,
		SRV: mongodb.SRVConfig{
			MaxHosts:    0,
			ServiceName: "mongodb",
		},
		ServerAPI: mongodb.ServerAPIConfig{
			Version:           "1",
			Strict:            false,
			DeprecationErrors: false,
		},
		ServerMonitoringMode:     "auto",
		DisableOCSPEndpointCheck: false,
		BSON: mongodb.BSONConfig{
			UseJSONStructTags:     false,
			ErrorOnInlineMap:      false,
			AllowTruncatingFloats: false,
		},
		AutoEncryption: mongodb.AutoEncryptionConfig{
			Enabled:              false,
			KMSProviders:         make(map[string]interface{}),
			SchemaMap:            make(map[string]interface{}),
			BypassAutoEncryption: false,
			ExtraOptions:         make(map[string]interface{}),
		},
	}

	actual := mongodb.DefaultConfig()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("DefaultConfig() = %+v, want %+v", actual, expected)
	}
}

func TestConfig_GetConnectTimeout(t *testing.T) {
	cfg := &mongodb.Config{ConnectTimeout: 5000}
	expected := 5000 * time.Millisecond
	if actual := cfg.GetConnectTimeout(); actual != expected {
		t.Errorf("GetConnectTimeout() = %v, want %v", actual, expected)
	}
}

func TestConfig_GetServerSelectionTimeout(t *testing.T) {
	cfg := &mongodb.Config{ServerSelectionTimeout: 10000}
	expected := 10000 * time.Millisecond
	if actual := cfg.GetServerSelectionTimeout(); actual != expected {
		t.Errorf("GetServerSelectionTimeout() = %v, want %v", actual, expected)
	}
}

func TestConfig_GetSocketTimeout(t *testing.T) {
	tests := []struct {
		name     string
		timeout  uint64
		expected time.Duration
	}{
		{"zero timeout", 0, 0},
		{"non-zero timeout", 15000, 15000 * time.Millisecond},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &mongodb.Config{SocketTimeout: tt.timeout}
			if actual := cfg.GetSocketTimeout(); actual != tt.expected {
				t.Errorf("GetSocketTimeout() = %v, want %v", actual, tt.expected)
			}
		})
	}
}

func TestConfig_GetHeartbeatInterval(t *testing.T) {
	cfg := &mongodb.Config{HeartbeatInterval: 2000}
	expected := 2000 * time.Millisecond
	if actual := cfg.GetHeartbeatInterval(); actual != expected {
		t.Errorf("GetHeartbeatInterval() = %v, want %v", actual, expected)
	}
}

func TestConfig_GetLocalThreshold(t *testing.T) {
	cfg := &mongodb.Config{LocalThreshold: 3000}
	expected := 3000 * time.Millisecond
	if actual := cfg.GetLocalThreshold(); actual != expected {
		t.Errorf("GetLocalThreshold() = %v, want %v", actual, expected)
	}
}

func TestConfig_GetTimeout(t *testing.T) {
	cfg := &mongodb.Config{Timeout: 25000}
	expected := 25000 * time.Millisecond
	if actual := cfg.GetTimeout(); actual != expected {
		t.Errorf("GetTimeout() = %v, want %v", actual, expected)
	}
}

func TestConfig_GetMaxConnIdleTime(t *testing.T) {
	cfg := &mongodb.Config{MaxConnIdleTime: 120000}
	expected := 120000 * time.Millisecond
	if actual := cfg.GetMaxConnIdleTime(); actual != expected {
		t.Errorf("GetMaxConnIdleTime() = %v, want %v", actual, expected)
	}
}

func TestConfig_GetWTimeout(t *testing.T) {
	cfg := &mongodb.Config{WriteConcern: mongodb.WriteConcernConfig{WTimeout: 7000}}
	expected := 7000 * time.Millisecond
	if actual := cfg.GetWTimeout(); actual != expected {
		t.Errorf("GetWTimeout() = %v, want %v", actual, expected)
	}
}
