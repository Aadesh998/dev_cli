package config

import (
	"cli/internal/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type LLMModel struct {
	ModelName string
	ModelAPI  string
}

var LlmClient LLMModel

func ConfigExists() bool {
	configPath := utils.GetDirConfigPath()
	configDir := filepath.Join(configPath, "dev_cli")
	configFile := filepath.Join(configDir, "config.yaml")

	_, err := os.Stat(configFile)
	return err == nil
}

// TODO: The error message for a missing model name is incorrect.
// It refers to 'claude_api' instead of 'model'. This should be corrected.
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

	LlmClient.ModelAPI = viper.GetString("apikey")
	if LlmClient.ModelAPI == "" {
		return fmt.Errorf("apikey is missing in config")
	}

	LlmClient.ModelName = viper.GetString("model")
	if LlmClient.ModelName == "" {
		return fmt.Errorf("model is missing in config")
	}

	return nil
}

// TODO: The function removes the config file before writing to it.
// This is not ideal as it can lead to data loss if the write fails.
// Consider using viper.WriteConfig() instead of viper.SafeWriteConfigAs()
// and handle the file creation and permissions more gracefully.
func SaveApiKey(model, key string) error {
	fmt.Printf("Model Name:= %s", model)
	configPath := utils.GetDirConfigPath()
	configDir := filepath.Join(configPath, "dev_cli")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	config_file := filepath.Join(configDir, "config.yaml")
	os.Remove(config_file)

	viper.Set("apikey", key)
	viper.Set("model", model)
	return viper.SafeWriteConfigAs(filepath.Join(configDir, "config.yaml"))
}
