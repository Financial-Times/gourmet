package main

import (
	"fmt"

	"github.com/Financial-Times/gourmet/config"
)

type appConfig struct {
	CustomConfig string `conf:"CUSTOM_CONFIG"`
	App          struct {
		SystemCode string `conf:"APP_SYSTEM_CODE" required:"true"`
		LogLevel   string `conf:"LOG_LEVEL" default:"INFO"`
	}
	Server struct {
		Enabed       bool `conf:"SERVER_ENABLED" default:"true"`
		Port         int  `conf:"SERVER_PORT" default:"8080"`
		ReadTimeout  int  `conf:"SERVER_READ_TIMEOUT" default:"10"`
		WriteTimeout int  `conf:"SERVER_WRITE_TIMETOUT" default:"15"`
		IdleTimeout  int  `conf:"SERVER_IDLE_TIMEOUT" default:"20"`
	}
}

func main() {
	confLoader := config.NewEnvConfigLoader()
	conf := appConfig{}
	err := confLoader.Load(&conf)
	fmt.Printf("%+v\n", err)
	fmt.Printf("%+v\n", conf)
}
