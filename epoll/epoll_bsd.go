// Package epoll lengdanran 2024/3/29 14:41
package epoll

import (
	"bufio"
	"fmt"
	"github.com/lengdanran/gredis/config"
	"github.com/lengdanran/gredis/redis/dbengine"
	"log"
	"log/slog"
	"os"
)

type BsdPollEventHandleFunc func(s *BsdSocket, dbEngine dbengine.DBEngine)

var DefaultBsdPollEventHandleFunc = func(s *BsdSocket, dbEngine dbengine.DBEngine) {
	reader := bufio.NewReader(s)
	for {
		line, err := reader.ReadString('\n')
		if line == "exit\n" {
			break
		}
		if err != nil {
			s.Close()
			break
		}
		log.Print("Read on socket ", s, "=>", line)
		s.Write([]byte(line))
	}
	// s.Close()
}

type BsdServer struct {
	Addr       string
	Port       int
	HandleFunc BsdPollEventHandleFunc
	dbEngine   dbengine.DBEngine
}

func NewBsdServer(config *config.GRedisServerConfig, handleFunc BsdPollEventHandleFunc, dbEngine dbengine.DBEngine) *BsdServer {
	svr := &BsdServer{
		Addr:       config.Addr,
		Port:       config.Port,
		dbEngine:   dbEngine,
		HandleFunc: handleFunc,
	}
	if svr.HandleFunc == nil {
		svr.HandleFunc = DefaultBsdPollEventHandleFunc
	}
	if svr.dbEngine == nil {
		svr.dbEngine = dbengine.NewRedisEngine()
	}
	return svr
}

func (svr *BsdServer) Start() {
	skt, err := Listen(svr.Addr, svr.Port)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create Socket: %v", err))
		os.Exit(1)
	}
	evtLoop, err := NewEventLoop(skt)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create kqueue: %v", err))
		os.Exit(1)
	}
	slog.Info("Server started. Waiting for incoming connections. ^C to exit.")
	evtLoop.Handle(svr.HandleFunc, svr.dbEngine)
}
