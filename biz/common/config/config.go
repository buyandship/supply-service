package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

type MercariConfig struct {
	ClientId     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	CallbackUrl  string `mapstructure:"callback_url"`
}

type MySQLConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Address  string `mapstructure:"addr"`
	DBName   string `mapstructure:"db_name"`
}

type RedisConfig struct {
	Address string `mapstructure:"addr"`
}

type OtelConfig struct {
	Endpoint    string `mapstructure:"endpoint"`
	Insecure    bool   `mapstructure:"insecure"`
	ServiceName string `mapstructure:"service_name"`
}

type ServerConfig struct {
	v          *viper.Viper
	Env        string        `mapstructure:"env"`
	Mysql      MySQLConfig   `mapstructure:"mysql"`
	Redis      RedisConfig   `mapstructure:"redis"`
	Mercari    MercariConfig `mapstructure:"mercari"`
	HmacSecret string        `mapstructure:"hmac_secret"`
	B4uToken   string        `mapstructure:"b4u_token"`
	Otel       OtelConfig    `mapstructure:"otel"`
}

var (
	GlobalServerConfig *ServerConfig
	once               sync.Once
)

func NewAppConfig(configPath string) (*ServerConfig, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("unable to read config: %w", err)
	}

	var config ServerConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to unmarshal config: %w", err)
	}

	config.v = v
	return &config, nil
}

// LoadGlobalConfig loads config as global configuration
func LoadGlobalConfig(configPath string) error {
	var err error
	once.Do(func() {
		GlobalServerConfig, err = NewAppConfig(configPath)
	})
	return err
}

func LoadTestConfig() {
	GlobalServerConfig = &ServerConfig{
		Env: "development",
		Mysql: MySQLConfig{
			Username: "root",
			Address:  "localhost:3306",
			DBName:   "mercari",
		},
		Redis: RedisConfig{
			Address: "localhost:6379",
		},
	}
}
