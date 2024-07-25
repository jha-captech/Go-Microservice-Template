package config

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfiguration(t *testing.T) {
	tests := map[string]struct {
		envVars       map[string]string
		expectedCfg   Configuration
		expectedError bool
	}{
		"success": {
			envVars: map[string]string{
				"ENV":                             "development",
				"LOG_LEVEL":                       "info",
				"DATABASE_NAME":                   "test_db",
				"DATABASE_USER":                   "test_user",
				"DATABASE_PASSWORD":               "test_password",
				"DATABASE_HOST":                   "localhost",
				"DATABASE_PORT":                   "5432",
				"DATABASE_RETRY_DURATION_SECONDS": "10",
			},
			expectedCfg: Configuration{
				Env:             "development",
				LogLevel:        slog.LevelInfo,
				DBName:          "test_db",
				DBUser:          "test_user",
				DBPassword:      "test_password",
				DBHost:          "localhost",
				DBPort:          "5432",
				DBRetryDuration: 10,
			},
			expectedError: false,
		},
		"missing required env": {
			envVars: map[string]string{
				"DATABASE_NAME":                   "test_db",
				"DATABASE_USER":                   "test_user",
				"DATABASE_PASSWORD":               "test_password",
				"DATABASE_HOST":                   "localhost",
				"DATABASE_PORT":                   "5432",
				"DATABASE_RETRY_DURATION_SECONDS": "10",
			},
			expectedCfg:   Configuration{},
			expectedError: true,
		},
		"invalid log level": {
			envVars: map[string]string{
				"ENV":                             "development",
				"LOG_LEVEL":                       "invalid_log_level",
				"DATABASE_NAME":                   "test_db",
				"DATABASE_USER":                   "test_user",
				"DATABASE_PASSWORD":               "test_password",
				"DATABASE_HOST":                   "localhost",
				"DATABASE_PORT":                   "5432",
				"DATABASE_RETRY_DURATION_SECONDS": "10",
			},
			expectedCfg:   Configuration{},
			expectedError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Set environment variables for the test
			for key, value := range tc.envVars {
				err := os.Setenv(key, value)
				assert.NoError(t, err)
			}

			// Cleanup environment variables after the test
			defer func() {
				for key := range tc.envVars {
					err := os.Unsetenv(key)
					assert.NoError(t, err)
				}
			}()

			cfg, err := New()
			if tc.expectedError {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedCfg, cfg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCfg, cfg)
			}
		})
	}
}
