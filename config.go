package irembopay

import (
	"fmt"
)

// EnvironmentType represents the IremboPay environment (sandbox or production)
type EnvironmentType string

const (
	// Sandbox is the testing environment
	Sandbox EnvironmentType = "sandbox"
	// Production is the live environment
	Production EnvironmentType = "production"
)

// Config holds the IremboPay API configuration
type Config struct {
	// Common configuration
	SecretKey   string          // Secret key for authentication
	APIVersion  string          // API version (default: "2")
	Environment EnvironmentType // Sandbox or Production
	Host        string          // API host URL
}

// NewConfig creates a new IremboPay configuration
func NewConfig(environment EnvironmentType, secretKey string, opts ...ConfigOption) (*Config, error) {
	// Default configuration based on environment
	config := &Config{
		Environment: environment,
		SecretKey:   secretKey,
		APIVersion:  "2", // Default API version
	}

	// Set environment-specific defaults
	switch environment {
	case Sandbox:
		config.Host = "api.sandbox.irembopay.com"
	case Production:
		config.Host = "api.irembopay.com"
	default:
		return nil, fmt.Errorf("invalid environment: %s", environment)
	}

	// Apply provided options
	for _, opt := range opts {
		opt(config)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// ConfigOption defines a function type for setting config options
type ConfigOption func(*Config)

// WithAPIVersion sets the API version
func WithAPIVersion(version string) ConfigOption {
	return func(c *Config) {
		c.APIVersion = version
	}
}

// WithHost sets the API host URL
func WithHost(host string) ConfigOption {
	return func(c *Config) {
		c.Host = host
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.SecretKey == "" {
		return fmt.Errorf("secret key is required")
	}
	if c.APIVersion == "" {
		return fmt.Errorf("API version is required")
	}
	if c.Host == "" {
		return fmt.Errorf("host is required")
	}
	return nil
}
