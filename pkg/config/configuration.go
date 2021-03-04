package config

import (
	"ratelimit/pkg/echorouter"
	"ratelimit/pkg/ratelimit"
	"ratelimit/pkg/redis"
	"ratelimit/pkg/zerolog"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var _config *Config

// Get get global config
func Get() *Config {
	return _config
}

// Set global config
func Set(config *Config) {
	_config = config
}

// Config ...
type Config struct {
	fx.Out
	Log       *zerolog.Config
	HTTP      *echorouter.Config
	Redis     redis.Config
	RateLimit ratelimit.Config `yaml:"ratelimit" mapstructure:"ratelimit"`
}

// New 讀取App 啟動程式設定檔
func New() (*Config, error) {
	viper.AutomaticEnv()

	configPath := viper.GetString("PROJ_DIR")
	if configPath == "" {
		configPath = "./"
	} else {
		configPath = configPath + "/configs"
	}

	configName := viper.GetString("CONFIG_NAME")
	if configName == "" {
		configName = "app"
	}

	viper.SetConfigName(configName)
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")

	var config Config

	if err := viper.ReadInConfig(); err != nil {
		log.Error().Msgf("Error reading config file, %s", err)
		return &config, err
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Error().Msgf("unable to decode into struct, %v", err)
		return &config, err
	}

	// 設定 rate limit interval 單位為 millisecond
	config.RateLimit.Interval *= time.Millisecond

	Set(&config)

	return _config, nil
}
