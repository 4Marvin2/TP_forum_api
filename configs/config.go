package configs

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type PostgresConfig struct {
	User     string
	Password string
	Port     string
	Host     string
	DBName   string
}

type TimeoutsConfig struct {
	WriteTimeout   time.Duration
	ReadTimeout    time.Duration
	ContextTimeout time.Duration
}

var (
	Postgres PostgresConfig
	Timeouts TimeoutsConfig
)

func SetConfig() {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	Postgres = PostgresConfig{
		Port:     viper.GetString(`postgres.port`),
		Host:     viper.GetString(`postgres.host`),
		User:     viper.GetString(`postgres.user`),
		Password: viper.GetString(`postgres.pass`),
		DBName:   viper.GetString(`postgres.name`),
	}

	Timeouts = TimeoutsConfig{
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
		ContextTimeout: time.Second * 2,
	}
}
