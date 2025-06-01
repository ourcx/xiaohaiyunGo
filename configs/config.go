package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `yaml:"app" mapstructure:"app"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Database DatabaseConfig `yaml:"database" mapstructure:"database"`
	Email    EmailConfig    `yaml:"email" mapstructure:"email"`
}

type AppConfig struct {
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	Driver string `yaml:"driver" mapstructure:"driver"`
	Source string `yaml:"source" mapstructure:"source"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	Db1      int    `mapstructure:"db1"`
	PoolSize int    `mapstructure:"pool_size"`
}

type EmailConfig struct {
	SmtpHost     string `mapstructure:"SMTP_HOST"`
	SmtpPort     int    `mapstructure:"SMTP_PORT"`
	SmtpUser     string `mapstructure:"SMTP_USER"`
	SmtpPassword string `mapstructure:"SMTP_PASSWORD"`
}

var Conf *Config

// LoadConfig 加载配置文件
func LoadConfig() error {

	// 设置配置文件路径和名称
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 将配置文件内容解析到 Conf 变量中
	Conf = &Config{}
	err = viper.Unmarshal(Conf)
	if err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	return nil
}
