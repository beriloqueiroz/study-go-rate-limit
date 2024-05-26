package config

import "github.com/spf13/viper"

type ConfigDefaultApiKeysLimitPerSecond struct {
	ApiKey         string `mapstructure:"api_key"`
	LimitPerSecond int    `mapstructure:"limit_per_second"`
}

type conf struct {
	WebServerPort                  string                             `mapstructure:"web_server_port"`
	DefaultLimitPerIPPerSecond     int                                `mapstructure:"default_limit_per_ip_per_second"`
	DefaultExpirationTimeInMinutes int                                `mapstructure:"default_expiration_time_in_minutes"`
	DefaultApiKeysLimitPerSecond   ConfigDefaultApiKeysLimitPerSecond `mapstructure:"default_api_keys_limit_per_second"`
}

func LoadConfig(paths []string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("json")
	for _, path := range paths {
		viper.AddConfigPath(path)
	}
	viper.SetConfigFile("config.json")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
