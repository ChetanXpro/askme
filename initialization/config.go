package initialization

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	OpenaiApiKeys  string
	HttpPort       int
	HttpProxy      string
	QdrantAddrGrpc string
	QdrantAddrHttp string
}

func LoadConfig(cfg string) *Config {
	viper.SetConfigFile(cfg)
	viper.ReadInConfig()
	viper.AutomaticEnv()

	config := &Config{
		OpenaiApiKeys:  getViperStringValue("OPENAI_KEY", ""),
		HttpPort:       getViperIntValue("HTTP_PORT", 8000),
		HttpProxy:      getViperStringValue("HTTP_PROXY", ""),
		QdrantAddrGrpc: getViperStringValue("QDRANT_ADDRGRPC", "localhost:6334"),
		QdrantAddrHttp: getViperStringValue("QDRANT_ADDRHTTP", "localhost:6333"),
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
