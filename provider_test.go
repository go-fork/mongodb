package mongodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.fork.vn/config/mocks"
	"go.fork.vn/di"
	diMocks "go.fork.vn/di/mocks"
)

// setupTestMongoConfig creates a MongoDB config for testing
func setupTestMongoConfig() *Config {
	return &Config{
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
		Auth: AuthConfig{
			Username:      "",
			Password:      "",
			AuthSource:    "",
			AuthMechanism: "",
		},
		TLS: TLSConfig{
			InsecureSkipVerify: false,
		},
		ReadPreference: ReadPreferenceConfig{
			Mode: "primary",
		},
		ReadConcern: ReadConcernConfig{
			Level: "",
		},
		WriteConcern: WriteConcernConfig{
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
	provider := NewServiceProvider()
	assert.NotNil(t, provider, "Expected service provider to be initialized")
}

func TestServiceProviderRegister(t *testing.T) {
	t.Run("registers mongodb services to container with config (mocked)", func(t *testing.T) {
		// Arrange
		mockContainer := diMocks.NewMockContainer(t)
		mockConfig := mocks.NewMockManager(t)
		mockApp := diMocks.NewMockApplication(t)
		mockApp.On("Container").Return(mockContainer).Maybe()

		testMongoConfig := setupTestMongoConfig()

		mockContainer.On("Bound", "config").Return(true).Maybe()
		mockContainer.On("Make", "config").Return(mockConfig, nil).Maybe()
		mockContainer.On("MustMake", "config").Return(mockConfig)
		mockConfig.EXPECT().UnmarshalKey("mongodb", mock.Anything).Run(func(_ string, out interface{}) {
			if cfg, ok := out.(*Config); ok {
				*cfg = *testMongoConfig
			}
		}).Return(nil)
		mockContainer.On("Singleton", "mongodb.manager", mock.Anything).Return(nil).Maybe()
		mockContainer.On("Singleton", "mongodb.client", mock.Anything).Return(nil).Maybe()
		mockContainer.On("Singleton", "mongodb.database", mock.Anything).Return(nil).Maybe()
		mockContainer.On("Instance", "mongodb", mock.Anything).Return(nil).Once()
		mockContainer.On("Instance", "mongodb.client", mock.Anything).Return(nil).Once()
		mockContainer.On("Instance", "mongodb.database", mock.Anything).Return(nil).Once()

		provider := NewServiceProvider()

		// Act
		initialProviders := provider.Providers()
		assert.Empty(t, initialProviders, "Expected 0 initial providers")
		provider.Register(mockApp)

		// Assert
		mockContainer.AssertCalled(t, "Instance", "mongodb", mock.Anything)
		mockContainer.AssertCalled(t, "Instance", "mongodb.client", mock.Anything)
		mockContainer.AssertCalled(t, "Instance", "mongodb.database", mock.Anything)
		finalProviders := provider.Providers()
		assert.Len(t, finalProviders, 3, "Expected 3 providers after registration")
	})

	t.Run("panics when config service is missing", func(t *testing.T) {
		mockContainer := diMocks.NewMockContainer(t)
		mockApp := diMocks.NewMockApplication(t)
		mockApp.On("Container").Return(mockContainer).Maybe()
		mockContainer.On("Bound", "config").Return(false).Maybe()
		mockContainer.On("MustMake", "config").Panic("config not bound").Maybe()
		mockContainer.On("Instance", mock.Anything, mock.Anything).Maybe()
		provider := NewServiceProvider()
		assert.Panics(t, func() {
			provider.Register(mockApp)
		}, "Expected provider.Register to panic when config is missing")
	})

	t.Run("panics when app doesn't have container", func(t *testing.T) {
		mockApp := diMocks.NewMockApplication(t)
		mockApp.On("Container").Return(nil).Maybe()
		provider := NewServiceProvider()
		assert.Panics(t, func() {
			provider.Register(mockApp)
		}, "Should panic when app doesn't have container")
	})

	t.Run("panics when config Make returns error", func(t *testing.T) {
		mockContainer := diMocks.NewMockContainer(t)
		mockApp := diMocks.NewMockApplication(t)
		mockApp.On("Container").Return(mockContainer).Maybe()
		mockContainer.On("Bound", "config").Return(true).Maybe()
		mockContainer.On("Make", "config").Return(nil, assert.AnError).Maybe()
		mockContainer.On("MustMake", "config").Panic("make error").Maybe()
		mockContainer.On("Instance", mock.Anything, mock.Anything).Maybe()
		provider := NewServiceProvider()
		assert.Panics(t, func() {
			provider.Register(mockApp)
		}, "Should panic when config Make returns error")
	})

	t.Run("panics when UnmarshalKey returns error", func(t *testing.T) {
		mockContainer := diMocks.NewMockContainer(t)
		mockConfig := mocks.NewMockManager(t)
		mockApp := diMocks.NewMockApplication(t)
		mockApp.On("Container").Return(mockContainer).Maybe()
		mockContainer.On("Bound", "config").Return(true).Maybe()
		mockContainer.On("Make", "config").Return(mockConfig, nil).Maybe()
		mockContainer.On("MustMake", "config").Return(mockConfig)
		mockConfig.EXPECT().UnmarshalKey("mongodb", mock.Anything).Return(assert.AnError)
		mockContainer.On("Instance", mock.Anything, mock.Anything).Maybe()
		provider := NewServiceProvider()
		assert.Panics(t, func() {
			provider.Register(mockApp)
		}, "Should panic when UnmarshalKey returns error")
	})
}

func TestServiceProviderBoot(t *testing.T) {
	t.Run("Boot doesn't panic", func(t *testing.T) {
		// Create DI container with config
		container := di.New()
		mockConfig := mocks.NewMockManager(t)

		// Setup expectations for UnmarshalKey
		mockConfig.EXPECT().UnmarshalKey("mongodb", mock.Anything).Run(func(_ string, out interface{}) {
			// Copy our test config to the output parameter
			if cfg, ok := out.(*Config); ok {
				*cfg = *setupTestMongoConfig()
			}
		}).Return(nil)

		container.Instance("config", mockConfig)

		// Create app and provider
		mockApp := diMocks.NewMockApplication(t)
		mockApp.On("Container").Return(container).Maybe()
		provider := NewServiceProvider()

		// First register the provider
		provider.Register(mockApp)

		// Then test that boot doesn't panic
		assert.NotPanics(t, func() {
			provider.Boot(mockApp)
		}, "Boot should not panic with valid configuration")
	})

	t.Run("Boot works without container", func(t *testing.T) {
		// Test with no container
		provider := NewServiceProvider()
		mockApp := diMocks.NewMockApplication(t)
		mockApp.On("Container").Return(nil).Maybe()

		// Should not panic
		assert.NotPanics(t, func() {
			provider.Boot(mockApp)
		}, "Boot should not panic with nil container")
	})
}

func TestServiceProviderBootWithNil(t *testing.T) {
	// Test Boot with nil app parameter
	provider := NewServiceProvider()

	// Should not panic with nil app
	assert.NotPanics(t, func() {
		provider.Boot(nil)
	}, "Boot should not panic with nil app parameter")
}

func TestServiceProviderProviders(t *testing.T) {
	// In the new implementation, providers are dynamically added during Register
	// So a freshly created provider should have an empty providers list
	provider := NewServiceProvider()
	providers := provider.Providers()

	assert.Empty(t, providers, "Expected empty providers list initially")

	// We test the dynamic registration of providers in TestServiceProviderRegister
}

func TestServiceProviderRequires(t *testing.T) {
	provider := NewServiceProvider()
	requires := provider.Requires()

	// MongoDB provider requires the config provider
	assert.Len(t, requires, 1, "Expected 1 required dependency")
	assert.Equal(t, "config", requires[0], "Expected required dependency to be 'config'")
}

func TestDynamicProvidersList(t *testing.T) {
	// This test verifies that providers are correctly registered in the dynamic list
	container := di.New()
	mockConfig := mocks.NewMockManager(t)

	// Setup expectations for UnmarshalKey
	mockConfig.EXPECT().UnmarshalKey("mongodb", mock.Anything).Run(func(_ string, out interface{}) {
		// Copy our test config to the output parameter
		if cfg, ok := out.(*Config); ok {
			*cfg = *setupTestMongoConfig()
		}
	}).Return(nil)

	container.Instance("config", mockConfig)
	mockApp := diMocks.NewMockApplication(t)
	mockApp.On("Container").Return(container).Maybe()
	provider := NewServiceProvider()

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

func TestServiceProviderInterfaceCompliance(t *testing.T) {
	// This test verifies that our concrete type implements the interface
	var _ ServiceProvider = (*serviceProvider)(nil)
	var _ di.ServiceProvider = (*serviceProvider)(nil)
}

func TestMockConfigManagerWithMongoConfig(t *testing.T) {
	// This test verifies that our mock config manager can be used with MongoDB config
	mockConfig := mocks.NewMockManager(t)
	testConfig := setupTestMongoConfig()

	// Setup expectations for the Has method
	mockConfig.EXPECT().Has("mongodb").Return(true)

	// Setup expectations for the Get method
	mockConfig.EXPECT().Get("mongodb").Return(testConfig, true)

	// Setup expectations for UnmarshalKey
	mockConfig.EXPECT().UnmarshalKey("mongodb", mock.Anything).Run(func(_ string, out interface{}) {
		// Copy our test config to the output parameter
		if cfg, ok := out.(*Config); ok {
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
	var outConfig Config
	err := mockConfig.UnmarshalKey("mongodb", &outConfig)
	assert.NoError(t, err, "UnmarshalKey should not return an error")

	// Verify our mock expectations were met
	mockConfig.AssertExpectations(t)
}

func TestServiceProviderWithMockery(t *testing.T) {
	t.Run("registers mongodb services using mockery", func(t *testing.T) {
		mockContainer := diMocks.NewMockContainer(t)
		mockConfig := mocks.NewMockManager(t)
		mockApp := diMocks.NewMockApplication(t)
		mockApp.On("Container").Return(mockContainer).Maybe()

		testMongoConfig := setupTestMongoConfig()
		mockConfig.EXPECT().UnmarshalKey("mongodb", mock.Anything).Run(func(_ string, out interface{}) {
			if cfg, ok := out.(*Config); ok {
				*cfg = *testMongoConfig
			}
		}).Return(nil)
		mockContainer.On("Bound", "config").Return(true).Maybe()
		mockContainer.On("Make", "config").Return(mockConfig, nil).Maybe()
		mockContainer.On("MustMake", "config").Return(mockConfig).Maybe()
		mockContainer.On("Instance", "mongodb", mock.Anything).Return().Maybe()
		mockContainer.On("Instance", "mongodb.manager", mock.Anything).Return().Maybe()
		mockContainer.On("Instance", "mongodb.client", mock.Anything).Return().Maybe()
		mockContainer.On("Instance", "mongodb.database", mock.Anything).Return().Maybe()
		mockContainer.On("Singleton", "mongodb.manager", mock.AnythingOfType("func(di.Container) interface{}")).Maybe()
		mockContainer.On("Singleton", "mongodb.client", mock.AnythingOfType("func(di.Container) interface{}")).Maybe()
		mockContainer.On("Singleton", "mongodb.database", mock.AnythingOfType("func(di.Container) interface{}")).Maybe()

		provider := NewServiceProvider()
		provider.Register(mockApp)
	})

	t.Run("handles registration error cases with mockery", func(t *testing.T) {
		mockContainer := diMocks.NewMockContainer(t)
		mockApp := diMocks.NewMockApplication(t)
		mockApp.On("Container").Return(mockContainer).Maybe()
		mockContainer.On("Bound", "config").Return(false).Maybe()
		mockContainer.On("MustMake", "config").Return(nil).Maybe()

		provider := NewServiceProvider()
		assert.Panics(t, func() {
			provider.Register(mockApp)
		}, "Should panic when config is not bound")

		mockContainer = diMocks.NewMockContainer(t)
		mockApp = diMocks.NewMockApplication(t)
		mockApp.On("Container").Return(mockContainer).Maybe()
		mockContainer.On("Bound", "config").Return(true).Maybe()
		mockContainer.On("Make", "config").Return(nil, assert.AnError).Maybe()
		mockContainer.On("MustMake", "config").Return(nil).Maybe()

		provider = NewServiceProvider()
		assert.Panics(t, func() {
			provider.Register(mockApp)
		}, "Should panic when config Make returns error")
	})
}

func TestBootWithMockery(t *testing.T) {
	// Test Boot with mockery
	mockContainer := diMocks.NewMockContainer(t)
	mockApp := diMocks.NewMockApplication(t)
	mockApp.On("Container").Return(mockContainer).Maybe()

	provider := NewServiceProvider()

	// Boot should not panic
	assert.NotPanics(t, func() {
		provider.Boot(mockApp)
	})
}
