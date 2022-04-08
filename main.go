package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/flexicon/nerfthis-tools/views"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	if err := ViperInit(); err != nil {
		return err
	}

	if err := CacheInit(); err != nil {
		return err
	}

	e := echo.New()
	e.Debug = viper.GetBool("debug")

	e.Use(middleware.Secure())
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "REQUEST: method=${method}, status=${status}, uri=${uri}, latency=${latency_human}\n",
	}))
	e.Renderer = views.NewRenderer()

	SetupRoutes(e.Group(""))

	return e.Start(fmt.Sprintf(":%d", viper.GetInt("port")))
}

// ViperInit loads environment variables and sets up needed config defaults.
func ViperInit() error {
	// Prepare for Environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("port", 80)
	viper.SetDefault("debug", false)
	viper.SetDefault("cache.ttl", 300) // 5 minutes in seconds

	return nil
}
