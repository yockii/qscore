package config

import (
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type config struct {
	*viper.Viper
}

var DefaultInstance = &config{viper.New()}

func init() {
	DefaultInstance.SetConfigName("config")
	DefaultInstance.AddConfigPath("./conf")
	if err := DefaultInstance.ReadInConfig(); err != nil {
		logger.Warnf("No config file: %s ", err)
	}
	DefaultInstance.AutomaticEnv()
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
func (c *config) GetUint64(key string) uint64 {
	return c.Viper.GetUint64(key)
}
func (c *config) GetFloat64(key string) float64 {
	return c.Viper.GetFloat64(key)
}
func (c *config) IsSet(key string) bool {
	return c.Viper.IsSet(key)
}
func (c *config) GetStringSlice(key string) []string {
	return c.Viper.GetStringSlice(key)
}
func (c *config) GetStringMapString(key string) map[string]string {
	return c.Viper.GetStringMapString(key)
}
func (c *config) GetDuration(key string) time.Duration {
	return c.Viper.GetDuration(key)
}
func (c *config) WatchConfig() {
	c.Viper.WatchConfig()
}

/////////////////////////////////////////////////////////////////////
//////// 默认配置获取 //////

func SetConfigName(name string) {
	DefaultInstance.SetConfigName(name)
}
func AddConfigPath(path string) {
	DefaultInstance.AddConfigPath(path)
}
func ReadInConfig() error {
	return DefaultInstance.ReadInConfig()
}
func AutomaticEnv() {
	DefaultInstance.AutomaticEnv()
}

func GetString(key string) string {
	return DefaultInstance.GetString(key)
}
func GetInt(key string) int {
	return DefaultInstance.GetInt(key)
}
func GetBool(key string) bool {
	return DefaultInstance.GetBool(key)
}
func GetUint(key string) uint {
	return DefaultInstance.GetUint(key)
}
func GetUint64(key string) uint64 {
	return DefaultInstance.GetUint64(key)
}
func GetFloat64(key string) float64 {
	return DefaultInstance.GetFloat64(key)
}
func IsSet(key string) bool {
	return DefaultInstance.IsSet(key)
}
func GetStringSlice(key string) []string {
	return DefaultInstance.GetStringSlice(key)
}
func GetIntSlice(key string) []int {
	return DefaultInstance.GetIntSlice(key)
}
func GetStringMapString(key string) map[string]string {
	return DefaultInstance.GetStringMapString(key)
}
func GetDuration(key string) time.Duration {
	return DefaultInstance.GetDuration(key)
}
func WatchConfig() {
	DefaultInstance.WatchConfig()
}
