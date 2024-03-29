//go:build linux
// +build linux

// Package epoll lengdanran 2024/3/29 16:50
package epoll

import (
	"bufio"
	"fmt"
	"github.com/lengdanran/gredis/config"
	"github.com/lengdanran/gredis/redis/dbengine"
	"github.com/lengdanran/gredis/redis/parser"
	"github.com/lengdanran/gredis/redis/protocol"
	"golang.org/x/sys/unix"
	"io"
	"log"
	"log/slog"
	"net"
	"os"
	"reflect"
	"strings"
	"sync"
	"syscall"
)

type epoll struct {
	fd          int
	connections map[int]net.Conn
	lock        *sync.RWMutex
}

type EpollServer struct {
	Addr     string
	Port     int
	dbEngine dbengine.DBEngine
}

func NewEpollServer(config *config.GRedisServerConfig, dbEngine dbengine.DBEngine) *EpollServer {
	svr := &EpollServer{
		Addr:     config.Addr,
		Port:     config.Port,
		dbEngine: dbEngine,
	}
	if svr.dbEngine == nil {
		svr.dbEngine = dbengine.NewRedisEngine()
	}
	return svr
}

func (svr *EpollServer) Start() {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", svr.Addr, svr.Port))
	if err != nil {
		log.Printf("Failed to create Socket: %v", err)
		os.Exit(1)
	}
	epoller, err := NewEpoll()
	go start(epoller, svr.dbEngine)
	for {
		conn, e := ln.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				log.Printf("accept temp err: %v", ne)
				continue
			}
			log.Printf("accept err: %v", e)
			return
		}
		if err := epoller.Add(conn); err != nil {
			log.Printf("failed to add connection %v", err)
			conn.Close()
		}
	}
}

func start(epoller *epoll, engine dbengine.DBEngine) {
	for {
		connections, err := epoller.Wait()
		if err != nil {
			// log.Printf("failed to epoll wait %v", err)
			continue
		}
		for _, conn := range connections {
			if conn == nil {
				break
			}
			reader := bufio.NewReader(conn)
			ch := make(chan *parser.Payload)
			go parser.ParseRedisRequestStream(reader, ch)
			for payload := range ch {
				if payload.Err != nil {
					if payload.Err == io.EOF ||
						payload.Err == io.ErrUnexpectedEOF ||
						strings.Contains(payload.Err.Error(), "use of closed network connection") {
						// 断开连接，终止处理
						slog.Info("connection closed!")
						return
					}
					errReply := protocol.MakeErrReply(payload.Err.Error())
					_, err := conn.Write(errReply.ToBytes())
					if err != nil {
						slog.Error(fmt.Sprintf("Write data to socket error: %v", err))
						return
					}
					continue
				} else {
					if payload.Data == nil {
						slog.Error("empty payload")
						continue
					}
					r, ok := payload.Data.(*protocol.MultiBulkReply)
					if !ok {
						slog.Error("require multi bulk protocol")
						continue
					}
					// 数据库处理引擎处理操作命令，并得到返回结果
					result := engine.Exec(r.Args)
					slog.Info(string(result.ToBytes()))
					if result != nil {
						_, _ = conn.Write(result.ToBytes())
					} else {
						_, _ = conn.Write([]byte("-ERR unknown\r\n"))
					}
				}
			}
		}
	}
}

func NewEpoll() (*epoll, error) {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &epoll{
		fd:          fd,
		lock:        &sync.RWMutex{},
		connections: make(map[int]net.Conn),
	}, nil
}

func (e *epoll) Add(conn net.Conn) error {
	// Extract file descriptor associated with the connection
	fd := socketFD(conn)
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: unix.POLLIN | unix.POLLHUP, Fd: int32(fd)})
	if err != nil {
		return err
	}
	e.lock.Lock()
	defer e.lock.Unlock()
	e.connections[fd] = conn
	if len(e.connections)%100 == 0 {
		log.Printf("Total number of connections: %v", len(e.connections))
	}
	return nil
}

func (e *epoll) Remove(conn net.Conn) error {
	fd := socketFD(conn)
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		return err
	}
	e.lock.Lock()
	defer e.lock.Unlock()
	delete(e.connections, fd)
	if len(e.connections)%100 == 0 {
		log.Printf("Total number of connections: %v", len(e.connections))
	}
	return nil
}

func (e *epoll) Wait() ([]net.Conn, error) {
	events := make([]unix.EpollEvent, 100)
	n, err := unix.EpollWait(e.fd, events, 100)
	if err != nil {
		return nil, err
	}
	e.lock.RLock()
	defer e.lock.RUnlock()
	var connections []net.Conn
	for i := 0; i < n; i++ {
		conn := e.connections[int(events[i].Fd)]
		connections = append(connections, conn)
	}
	return connections, nil
}

func socketFD(conn net.Conn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}
