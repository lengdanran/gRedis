// Package main lengdanran 2024/3/27 16:53
//go:build darwin

package main

import (
	"github.com/lengdanran/gredis/config"
	"github.com/lengdanran/gredis/epoll"
	"github.com/lengdanran/gredis/redis/handler"
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
	// for darwin
	server := epoll.NewBsdServer(config.ServerConfig, handler.RedisBsdPollEventHandleFunc, nil)
	server.Start()
}
