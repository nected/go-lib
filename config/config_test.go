package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAppConfig(t *testing.T) {
	// Call the function under test
	config := NewAppConfig()

	// Verify the expected values
	expectedConfigPaths := defaultConfigPaths
	assert.Equal(t, expectedConfigPaths, config.configPaths, "NewAppConfig() returned configPaths %v, expected %v", config.configPaths, expectedConfigPaths)

	// Set and verify app name
	expectedAppName := "myapp"
	config.SetAppName(expectedAppName)
	assert.Equal(t, expectedAppName, config.GetAppName(), "SetAppName() and GetAppName() returned %s, expected %s", config.GetAppName(), expectedAppName)
	assert.Equal(t, expectedAppName, config.GetConfigFileName(), "SetAppName() and GetConfigFileName() returned %s, expected %s", config.GetConfigFileName(), expectedAppName)

	// Set and verify config file path
	expectedConfigFilePath := "/etc/myapp"
	config.SetConfigFilePath(expectedConfigFilePath)
	assert.Equal(t, expectedConfigFilePath, config.GetConfigFilePath(), "SetConfigFilePath() and GetConfigFilePath() returned %s, expected %s", config.GetConfigFilePath(), expectedConfigFilePath)

	// Set and verify config file name
	expectedConfigFileName := "myappfile"
	config.SetConfigFileName(expectedConfigFileName)
	assert.Equal(t, expectedConfigFileName, config.GetConfigFileName(), "SetConfigFileName() and GetConfigFileName() returned %s, expected %s", config.GetConfigFileName(), expectedConfigFileName)
	assert.NotEqual(t, config.GetAppName(), config.GetConfigFileName(), "SetConfigFileName() and GetConfigFileName() returned %s, expected %s", config.GetConfigFileName(), expectedAppName)
}
