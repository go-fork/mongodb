package mongodb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	config_mocks "go.fork.vn/config/mocks"
	di_mocks "go.fork.vn/di/mocks"
	"go.fork.vn/mongodb"
)

// setupTestMongoConfig creates a MongoDB config for testing
func setupTestMongoConfig() *mongodb.Config {
	return &mongodb.Config{
		URI:                    "mongodb://localhost:27017",
		Database:               "testdb",
		ConnectTimeout:         10000,
		MaxPoolSize:            10,
		MinPoolSize:            1,
		MaxConnIdleTime:        300000,
		HeartbeatInterval:      30000,
		ServerSelectionTimeout: 30000,
		SocketTimeout:          5000,
		LocalThreshold:         15000,
		Auth: mongodb.AuthConfig{
			Username:      "",
			Password:      "",
			AuthSource:    "",
			AuthMechanism: "",
		},
		TLS: mongodb.TLSConfig{
			InsecureSkipVerify: false,
		},
		ReadPreference: mongodb.ReadPreferenceConfig{
			Mode: "primary",
		},
		ReadConcern: mongodb.ReadConcernConfig{
			Level: "",
		},
		WriteConcern: mongodb.WriteConcernConfig{
			W:        1,
			WTimeout: 0,
			Journal:  false,
		},
		AppName:      "",
		Direct:       false,
		ReplicaSet:   "",
		Compressors:  []string{},
		RetryWrites:  true,
		RetryReads:   true,
		LoadBalanced: false,
	}
}

func TestNewServiceProvider(t *testing.T) {
	t.Run("creates real service provider", func(t *testing.T) {
		provider := mongodb.NewServiceProvider()
		assert.NotNil(t, provider, "Expected service provider to be initialized")
	})
}

func TestServiceProvider_Register(t *testing.T) {
	t.Run("registers mongodb services to container with config", func(t *testing.T) {
		// Arrange
		mockContainer := di_mocks.NewMockContainer(t)
		mockConfig := config_mocks.NewMockManager(t)
		mockApp := di_mocks.NewMockApplication(t)
		mockApp.On("Container").Return(mockContainer)

		testMongoConfig := setupTestMongoConfig()

		mockContainer.On("MustMake", "config").Return(mockConfig)
		mockConfig.EXPECT().UnmarshalKey("mongodb", mock.Anything).Run(func(_ string, out interface{}) {
			if cfg, ok := out.(*mongodb.Config); ok {
				*cfg = *testMongoConfig
			}
		}).Return(nil)
		mockContainer.On("Instance", "mongodb", mock.Anything).Return(nil)
		mockContainer.On("Instance", "mongodb.client", mock.Anything).Return(nil)
		mockContainer.On("Instance", "mongodb.database", mock.Anything).Return(nil)

		provider := mongodb.NewServiceProvider()

		// Act & Assert
		assert.NotPanics(t, func() {
			provider.Register(mockApp)
		}, "Register should not panic with valid configuration")

		// Verify that the container methods were called
		mockContainer.AssertCalled(t, "Instance", "mongodb", mock.Anything)
		mockContainer.AssertCalled(t, "Instance", "mongodb.client", mock.Anything)
		mockContainer.AssertCalled(t, "Instance", "mongodb.database", mock.Anything)
	})

	t.Run("panics when config service is missing", func(t *testing.T) {
		mockContainer := di_mocks.NewMockContainer(t)
		mockApp := di_mocks.NewMockApplication(t)
		mockApp.On("Container").Return(mockContainer)
		mockContainer.On("MustMake", "config").Panic("config not bound")

		provider := mongodb.NewServiceProvider()

		assert.Panics(t, func() {
			provider.Register(mockApp)
		}, "Expected provider.Register to panic when config is missing")
	})

	t.Run("panics when app doesn't have container", func(t *testing.T) {
		mockApp := di_mocks.NewMockApplication(t)
		mockApp.On("Container").Return(nil)

		provider := mongodb.NewServiceProvider()

		assert.Panics(t, func() {
			provider.Register(mockApp)
		}, "Should panic when app doesn't have container")
	})

	t.Run("panics when UnmarshalKey returns error", func(t *testing.T) {
		mockContainer := di_mocks.NewMockContainer(t)
		mockConfig := config_mocks.NewMockManager(t)
		mockApp := di_mocks.NewMockApplication(t)
		mockApp.On("Container").Return(mockContainer)
		mockContainer.On("MustMake", "config").Return(mockConfig)
		mockConfig.EXPECT().UnmarshalKey("mongodb", mock.Anything).Return(assert.AnError)

		provider := mongodb.NewServiceProvider()

		assert.Panics(t, func() {
			provider.Register(mockApp)
		}, "Should panic when UnmarshalKey returns error")
	})
}

func TestServiceProvider_Boot(t *testing.T) {
	t.Run("Boot doesn't panic with valid app", func(t *testing.T) {
		provider := mongodb.NewServiceProvider()
		mockApp := di_mocks.NewMockApplication(t)

		// Boot should not panic with valid app
		assert.NotPanics(t, func() {
			provider.Boot(mockApp)
		}, "Boot should not panic with valid app")
	})

	t.Run("Boot panics with nil app", func(t *testing.T) {
		provider := mongodb.NewServiceProvider()

		// Should panic with nil app
		assert.Panics(t, func() {
			provider.Boot(nil)
		}, "Boot should panic with nil app parameter")
	})
}

func TestServiceProvider_BootWithNil(t *testing.T) {
	// Test Boot with nil app parameter
	provider := mongodb.NewServiceProvider()

	// Should panic with nil app
	assert.Panics(t, func() {
		provider.Boot(nil)
	}, "Boot should panic with nil app parameter")
}

func TestServiceProvider_Providers(t *testing.T) {
	// Test that providers list is initially empty and populated after registration
	provider := mongodb.NewServiceProvider()

	providers := provider.Providers()
	assert.Empty(t, providers, "Expected empty providers list initially")
}

func TestServiceProvider_Requires(t *testing.T) {
	provider := mongodb.NewServiceProvider()

	requires := provider.Requires()

	// MongoDB provider requires the config provider
	assert.Len(t, requires, 1, "Expected 1 required dependency")
	assert.Equal(t, "config", requires[0], "Expected required dependency to be 'config'")
}

func TestDynamicProvidersList(t *testing.T) {
	// Test that providers are correctly registered in the dynamic list
	mockContainer := di_mocks.NewMockContainer(t)
	mockConfig := config_mocks.NewMockManager(t)
	mockApp := di_mocks.NewMockApplication(t)
	mockApp.On("Container").Return(mockContainer)

	testMongoConfig := setupTestMongoConfig()
	mockContainer.On("MustMake", "config").Return(mockConfig)
	mockConfig.EXPECT().UnmarshalKey("mongodb", mock.Anything).Run(func(_ string, out interface{}) {
		if cfg, ok := out.(*mongodb.Config); ok {
			*cfg = *testMongoConfig
		}
	}).Return(nil)
	mockContainer.On("Instance", "mongodb", mock.Anything).Return(nil)
	mockContainer.On("Instance", "mongodb.client", mock.Anything).Return(nil)
	mockContainer.On("Instance", "mongodb.database", mock.Anything).Return(nil)

	provider := mongodb.NewServiceProvider()

	// Initially empty providers list
	initialProviders := provider.Providers()
	assert.Empty(t, initialProviders, "Expected 0 initial providers")

	// Register provider
	provider.Register(mockApp)

	// Check providers list after registration
	providers := provider.Providers()

	// We expect 3 entries: mongodb, mongodb.client, mongodb.database
	expectedItems := []string{"mongodb", "mongodb.client", "mongodb.database"}
	for _, expected := range expectedItems {
		assert.Contains(t, providers, expected, "Expected to find '%s' in providers list", expected)
	}

	// Length should match too
	assert.Len(t, providers, len(expectedItems), "Expected %d providers", len(expectedItems))
}

func TestMockConfigManagerWithMongoConfig(t *testing.T) {
	// This test verifies that our mock config manager can be used with MongoDB config
	mockConfig := config_mocks.NewMockManager(t)
	testConfig := setupTestMongoConfig()

	// Setup expectations for the Has method
	mockConfig.EXPECT().Has("mongodb").Return(true)

	// Setup expectations for the Get method
	mockConfig.EXPECT().Get("mongodb").Return(testConfig, true)

	// Setup expectations for UnmarshalKey
	mockConfig.EXPECT().UnmarshalKey("mongodb", mock.Anything).Run(func(_ string, out interface{}) {
		// Copy our test config to the output parameter
		if cfg, ok := out.(*mongodb.Config); ok {
			*cfg = *testConfig
		}
	}).Return(nil)

	// Test Has method
	assert.True(t, mockConfig.Has("mongodb"), "Has should return true for mongodb key")

	// Test Get method
	value, exists := mockConfig.Get("mongodb")
	assert.True(t, exists, "Should find the mongodb key")
	assert.Equal(t, testConfig, value, "Should return our test config")

	// Test UnmarshalKey method
	var outConfig mongodb.Config
	err := mockConfig.UnmarshalKey("mongodb", &outConfig)
	assert.NoError(t, err, "UnmarshalKey should not return an error")

	// Verify our mock expectations were met
	mockConfig.AssertExpectations(t)
}
