package main

import (
	"fmt"
	"log"

	"github.com/pimentafm/fc-ratelimiter/config"
	"github.com/redis/go-redis/v9"
)

func main() {
	var cfg config.Config
	viperCfg := config.NewViper("env")
	viperCfg.ReadViper(&cfg)

	newWebServer := CreateWebServer(
		&cfg,
		redis.NewClient(
			&redis.Options{
				Addr: fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
				DB:   cfg.Redis.Db,
			},
		),
	)

	log.Println("Starting web server on port", cfg.App.Port)
	newWebServer.Start()
}
