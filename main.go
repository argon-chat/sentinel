package main

import (
	"fmt"
	"github.com/argon-chat/sentinel/pkg/config"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("settings")
	viper.SetConfigType("json")
	viper.AddConfigPath("/etc/sentinel/")
	viper.AddConfigPath("$HOME/.sentinel")
	viper.AddConfigPath(".")
	var defaultPort float64 = 3000
	viper.SetDefault("server.port", defaultPort)
	viper.SetDefault("server.route", "/tunnel")
	viper.SetDefault("projects", map[string]string{})
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	config.Parse(viper.Get("projects"), viper.Get("server"))
}

func main() {
	fmt.Println(config.Instance)
}
