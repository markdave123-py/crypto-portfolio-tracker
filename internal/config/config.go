package config

import (
	"github.com/caarlos0/env/v10"
)

type Config struct {
	App  AppConfig
	HTTP HTTPConfig
	Log  LogConfig

	CoinGecko CoinGeckoConfig
	Redis     RedisConfig
	DB        DBConfig
	Pricing   PricingConfig
	EtherScan EtherScanConfig
}

type AppConfig struct {
	Name string `env:"APP_NAME" envDefault:"crypto-portfolio"`
	Env  string `env:"APP_ENV" envDefault:"local"`
}

// type RetryConfig struct {
// 	MaxRetries int `env:"MAX_RETRIES" envDefault:"3"`
// 	BaseDelay  int `env:"BASE_DELAY" envDefault:"500"`
// 	MaxDelay   int `env:"MAX_DELAY" envDefault:"800"`
// }

type HTTPConfig struct {
	Port string `env:"HTTP_PORT" envDefault:"8080"`
}

type LogConfig struct {
	Level string `env:"LOG_LEVEL" envDefault:"info"`
}

type CoinGeckoConfig struct {
	APIKey  string `env:"COINGECKO_API_KEY,required"`
	BaseURL string `env:"COINGECKO_BASE_URL" envDefault:"https://api.coingecko.com/api/v3"`
}

type RedisConfig struct {
	URL string `env:"REDIS_URL" envDefault:"redis://localhost:6379/0"`
}

type PricingConfig struct {
	CacheTTLSeconds int `env:"CACHE_TTL_SECONDS" envDefault:"30"`
}

type EtherScanConfig struct {
	APIKey  string `env:"ETHERSCAN_API_KEY,required"`
	BaseURL string `env:"ETHERSCAN_BASE_URL" envDefault:"https://api.etherscan.io/v2/api"`
}

type DBConfig struct {
	URL string `env:"DATABASE_URL"`
}

func (dbCfg *DBConfig) ConnectionString() string {
	return dbCfg.URL
}

// func (retryCfg *RetryConfig) GetConfig() *RetryConfig {
// 	return retryCfg
// }

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
