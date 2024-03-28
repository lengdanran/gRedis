// Package main lengdanran 2024/3/27 16:53
package main

import (
	"github.com/lengdanran/gredis/config"
	"github.com/lengdanran/gredis/redis/handler"
	"github.com/lengdanran/gredis/server"
	"log/slog"
)

const BANNER = `
       ____          _ _     
  __ _|  _ \ ___  __| (_)___ 
 / _' | |_) / _ \/ _' | / __|
| (_| |  _ <  __/ (_| | \__ \
 \__, |_| \_\___|\__,_|_|___/
 |___/                       
`

func main() {
	slog.Info(BANNER)
	// read configuration
	slog.Info(config.ServerConfig.RunId)
	// start gredis server
	err := server.StartServer(config.ServerConfig.Addr, config.ServerConfig.Port, handler.NewRedisHandler())
	if err != nil {
		return
	}
}
