package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var defaultConfigPaths = []string{
	"/etc/{{.appName}}",
	"$HOME/.{{.appName}}",
}

type appConfig struct {
	appName     string
	filePath    string
	fileName    string
	envPrefix   string
	configPaths []string
}

func NewAppConfig() *appConfig {
	setCWDToConfigPath()
	return &appConfig{
		configPaths: defaultConfigPaths,
	}
}

func (c *appConfig) GetConfigFilePath() string {
	return c.filePath
}

func (c *appConfig) SetConfigFilePath(configFilePath string) *appConfig {
	c.filePath = configFilePath
	return c
}

func (c *appConfig) GetConfigFileName() string {
	return c.fileName
}

func (c *appConfig) SetConfigFileName(configFileName string) *appConfig {
	c.fileName = configFileName
	return c
}

func (c *appConfig) GetEnvPrefix() string {
	return c.envPrefix
}

func (c *appConfig) SetEnvPrefix(envPrefix string) *appConfig {
	c.envPrefix = envPrefix
	return c
}

func (c *appConfig) GetAppName() string {
	return c.appName
}

func (c *appConfig) SetAppName(appName string) *appConfig {
	c.appName = appName
	if c.fileName == "" {
		c.fileName = appName
	}
	return c
}

func (c *appConfig) GetConfigFile() string {
	if c.filePath == "" || c.fileName == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s", c.filePath, c.fileName)
}

func (c *appConfig) GetConfigPaths() []string {
	return c.configPaths
}

func (c *appConfig) SetConfigPaths(configPaths []string) *appConfig {
	c.configPaths = append(c.configPaths, configPaths...)
	return c
}

func (c *appConfig) LoadConfig() {
	// Load config from file using viper
	if c.GetConfigFile() != "" {
		viper.SetConfigFile(c.GetConfigFile())
	}
}

func setCWDToConfigPath() {
	if cwd, err := os.Getwd(); err != nil {
		fmt.Println("Error getting current working directory:", err)
	} else {
		defaultConfigPaths = append(defaultConfigPaths, cwd)
	}
}
