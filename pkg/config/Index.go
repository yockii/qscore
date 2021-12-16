package config

import (
	"github.com/spf13/viper"

	"github.com/yockii/qscore/pkg/logger"
)

type config struct {
	*viper.Viper
}

var defaultConfig = &config{viper.New()}

func init() {
	defaultConfig.SetConfigName("config")
	defaultConfig.AddConfigPath("./conf")
	if err := defaultConfig.ReadInConfig(); err != nil {
		logger.Warnf("No config file: %s ", err)
	}
	defaultConfig.AutomaticEnv()
}

func (c *config) SetConfigName(name string) {
	c.Viper.SetConfigName(name)
}
func (c *config) AddConfigPath(path string) {
	c.Viper.AddConfigPath(path)
}
func (c *config) ReadInConfig() error {
	return c.Viper.ReadInConfig()
}
func (c *config) AutomaticEnv() {
	c.Viper.AutomaticEnv()
}

func (c *config) GetString(key string) string {
	return c.Viper.GetString(key)
}
func (c *config) GetInt(key string) int {
	return c.Viper.GetInt(key)
}
func (c *config) GetBool(key string) bool {
	return c.Viper.GetBool(key)
}
func (c *config) GetUint(key string) uint {
	return c.Viper.GetUint(key)
}
func (c *config) IsSet(key string) bool {
	return c.Viper.IsSet(key)
}

/////////////////////////////////////////////////////////////////////
//////// 默认配置获取 //////

func SetConfigName(name string) {
	defaultConfig.SetConfigName(name)
}
func AddConfigPath(path string) {
	defaultConfig.AddConfigPath(path)
}
func ReadInConfig() error {
	return defaultConfig.ReadInConfig()
}
func AutomaticEnv() {
	defaultConfig.AutomaticEnv()
}

func GetString(key string) string {
	return defaultConfig.GetString(key)
}
func GetInt(key string) int {
	return defaultConfig.GetInt(key)
}
func GetBool(key string) bool {
	return defaultConfig.GetBool(key)
}
func GetUint(key string) uint {
	return defaultConfig.GetUint(key)
}
func IsSet(key string) bool {
	return defaultConfig.IsSet(key)
}
