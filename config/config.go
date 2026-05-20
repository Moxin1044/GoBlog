package config

import "github.com/spf13/viper"

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	SMTP     SMTPConfig     `mapstructure:"smtp"`
	Feishu   FeishuConfig   `mapstructure:"feishu"`
	Upload   UploadConfig   `mapstructure:"upload"`
	AI       AIConfig       `mapstructure:"ai"`
	Monitor  MonitorConfig  `mapstructure:"monitor"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug/release
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"` // hours
	Issuer string `mapstructure:"issuer"`
}

type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

type FeishuConfig struct {
	Token   string `mapstructure:"token"`
	Enabled bool   `mapstructure:"enabled"`
}

type UploadConfig struct {
	Path       string   `mapstructure:"path"`
	MaxSize    int64    `mapstructure:"max_size"` // MB
	AllowTypes []string `mapstructure:"allow_types"`
}

type AIConfig struct {
	Enabled         bool `mapstructure:"enabled"`
	MaxDailyPerUser int  `mapstructure:"max_daily_per_user"`
}

type MonitorConfig struct {
	Enabled    bool `mapstructure:"enabled"`
	MapEnabled bool `mapstructure:"map_enabled"`
}

var AppConfig *Config

func InitConfig(path string) error {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	AppConfig = &Config{}
	return viper.Unmarshal(AppConfig)
}
