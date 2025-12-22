package config

import (
	"cli/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var LLMModel string

func ConfigExists() bool {
	configPath := utils.GetDirConfigPath()
	configDir := filepath.Join(configPath, "dev_cli")
	configFile := filepath.Join(configDir, "config.yaml")

	_, err := os.Stat(configFile)
	return err == nil
}

func LoadConfig() error {
	configPath := utils.GetDirConfigPath()
	configDir := filepath.Join(configPath, "dev_cli")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	LLMModel = viper.GetString("claude_api")
	if LLMModel == "" {
		return fmt.Errorf("claude_api is missing in config")
	}

	return nil
}

func SaveClaudeKey(key string) error {
	configPath := utils.GetDirConfigPath()
	configDir := filepath.Join(configPath, "dev_cli")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	viper.Set("claude_api", key)
	return viper.SafeWriteConfigAs(filepath.Join(configDir, "config.yaml"))
}
