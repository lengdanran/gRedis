// Package server lengdanran 2024/3/27 17:25
package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// ClientCounter 存放当前的客户端连接数量，可做最大连接的限制
var ClientCounter int32

// Handler 定义一个接口，用于处理客户端连接
type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}

func StartServer(addr string, port int, handler Handler) error {
	closeChan := make(chan bool)
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	// start a goroutine to handle the stop signals
	go func() {
		sig := <-sigCh
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- true
		}
	}()
	bindAddr := fmt.Sprintf("%s:%d", addr, port)
	listener, err := net.Listen("tcp", bindAddr)
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("bind: %s, start listening...", bindAddr))
	ListenAndServe(listener, handler, closeChan)
	return nil
}

func ListenAndServe(listener net.Listener, handler Handler, closeChan chan bool) {
	// listen signal
	errCh := make(chan error, 1)
	defer close(errCh)
	go func() {
		select {
		case <-closeChan:
			slog.Info("get exit signal")
		case er := <-errCh:
			slog.Info(fmt.Sprintf("accept error: %s", er.Error()))
		}
		slog.Info("shutting down...")
		_ = listener.Close() // listener.Accept() will return err immediately
		_ = handler.Close()
	}()

	ctx := context.Background()
	var waitDone sync.WaitGroup
	for {
		// accept a connection from client, if got the close signal, listener.Accept() will return err immediately
		conn, err := listener.Accept()
		if err != nil {
			// learn from net/http/serve.go#Serve()
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				slog.Info(fmt.Sprintf("accept occurs temporary error: %v, retry in 5ms", err))
				time.Sleep(5 * time.Millisecond)
				continue
			}
			errCh <- err
			break
		}
		// handle
		slog.Info("accept a new connection")
		ClientCounter++
		waitDone.Add(1)
		// 将连接请求的处理交给单独的goroutine中，主goroutine继续处理下一个请求
		go func() {
			defer func() {
				waitDone.Done()
				atomic.AddInt32(&ClientCounter, -1)
			}()
			handler.Handle(ctx, conn)
		}()
	}
	waitDone.Wait()
}
