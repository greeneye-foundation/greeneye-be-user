package config

import (
	"fmt"
	"os"

	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host        string `mapstructure:"host"`
		Port        string `mapstructure:"port"`
		Environment string `mapstructure:"environment"`
	} `mapstructure:"server"`

	MongoDB struct {
		URI           string `mapstructure:"uri"`
		Database      string `mapstructure:"database"`
		AuthMechanism string `mapstructure:"auth_mechanism"`
	} `mapstructure:"mongodb"`

	Redis struct {
		URI string `mapstructure:"uri"`
	} `mapstructure:"redis"`

	JWT struct {
		Secret     string `mapstructure:"secret"`
		Expiration int    `mapstructure:"expiration"`
	} `mapstructure:"jwt"`
}

func LoadConfig(configPath string) (*Config, error) {
	// Load .env files first
	loadDotEnvFiles()

	// Apply env variable replacement to config path
	configPath = replaceEnvVariables(configPath)

	// Viper configuration
	viper.SetConfigFile(configPath)

	// Replace environment variables in configuration values
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configPath)

	// Read and process configuration
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %v", err)
	}

	// Pre-process configuration values
	processConfigValues(viper.GetViper())

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %v", err)
	}
	return &config, nil
}

// Advanced configuration value processing
func processConfigValues(v *viper.Viper) {
	// Replace environment variables in specific configuration paths
	configPaths := []string{
		"mongodb.uri",
		"redis.uri",
		"jwt.secret",
	}

	for _, path := range configPaths {
		if value := v.GetString(path); value != "" {
			processedValue := replaceEnvVariables(value)
			v.Set(path, processedValue)
		}
	}
}

// Advanced env variable replacement
func replaceEnvVariables(input string) string {
	return os.Expand(input, func(v string) string {
		// Custom replacement logic
		if val, exists := os.LookupEnv(v); exists {
			return val
		}
		return fmt.Sprintf("${%s}", v)
	})
}

// Helper function to get env with default
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// loadDotEnvFiles loads multiple .env files
func loadDotEnvFiles() error {
	// Potential .env file locations
	envFiles := []string{
		".env",
		".env.local",
		".env.development",
		".env.production",
	}

	var errs []string
	for _, file := range envFiles {
		// Silently skip if file doesn't exist
		if err := godotenv.Load(file); err != nil && !os.IsNotExist(err) {
			errs = append(errs, fmt.Sprintf("Error loading %s: %v", file, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("env loading errors: %s", strings.Join(errs, "; "))
	}
	return nil
}
