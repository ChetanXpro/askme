package initialization

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	OpenaiApiKeys       string
	PineconeProjectName string
	PineconeIndexName   string
	PineconeEnvironment string
	PineconeAPIKEY      string
}

func LoadConfig(cfg string) *Config {
	viper.SetConfigFile(cfg)
	viper.ReadInConfig()
	viper.AutomaticEnv()

	config := &Config{
		OpenaiApiKeys:       getViperStringValue("openai_key", ""),
		PineconeProjectName: getViperStringValue("pinecone_project_name", ""),
		PineconeIndexName:   getViperStringValue("pinecone_index_name", ""),
		PineconeEnvironment: getViperStringValue("pinecone_environment", ""),
		PineconeAPIKEY:      getViperStringValue("pinecone_api_key", ""),
	}

	return config
}

func getViperStringValue(key string, defaultValue string) string {
	value := viper.GetString(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getViperIntValue(key string, defaultValue int) int {
	value := viper.GetString(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		fmt.Printf("Invalid value for %s, using default value %d\n", key, defaultValue)
		return defaultValue
	}
	return intValue
}
