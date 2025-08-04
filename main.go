package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("settings")
	viper.SetConfigType("json")
	viper.AddConfigPath("/etc/sentinel/")
	viper.AddConfigPath("$HOME/.sentinel")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func main() {
}
