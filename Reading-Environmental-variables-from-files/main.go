package main

import (
	"log"
	"github.com/spf13/viper"
)

type Config struct {
	DBUsername string `mapstructure:"DB_USERNAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
}

func main() {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("could not read env file:", err)
		return
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Println("error unmarshalling configuration variables:", err)
		return
	}
	log.Println("successfully read from env file:", config)
}

