package config

import "os"

type Config struct {
	AppEnv     string
	ServerPort string

	Settings *Settings
}

func InitConfig() *Config {
	instance := new(Config)

	instance.AppEnv = os.Getenv("APP_ENV")
	instance.ServerPort = os.Getenv("SERVER_PORT")

	return instance
}
