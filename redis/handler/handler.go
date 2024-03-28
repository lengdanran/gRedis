// Package handler lengdanran 2024/3/27 18:43
package handler

import (
	"context"
	"github.com/lengdanran/gredis/lib/sync/atomic"
	"github.com/lengdanran/gredis/redis/connection"
	"github.com/lengdanran/gredis/redis/dbengine"
	"github.com/lengdanran/gredis/redis/parser"
	"github.com/lengdanran/gredis/redis/protocol"
	"github.com/lengdanran/gredis/server"
	"io"
	"log/slog"
	"net"
	"strings"
	"sync"
)

type Handler struct {
	server.Handler
	activeConn sync.Map          // *client -> placeholder
	dbEngine   dbengine.DBEngine // 数据库引擎
	closing    atomic.Boolean    // 关闭客户端,拒绝新客户端和新请求
}

// NewRedisHandler creates a Handler instance
func NewRedisHandler() *Handler {
	return &Handler{
		dbEngine: &dbengine.SimpleEngine{},
	}
}

func (h *Handler) closeClient(client *connection.Connection) {
	_ = client.Close()
	h.dbEngine.AfterClientClose(client)
	h.activeConn.Delete(client)
}

// Handle 接受一个连接，处理请求
func (h *Handler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		// 拒绝新客户端和新请求
		_ = conn.Close()
		return
	}
	// 客户端封装
	client := connection.NewConn(conn)
	// 存放客户端到activeConn
	h.activeConn.Store(client, struct{}{})
	// 解析redis客户端的请求，ch里面会返回该连接按照每行解析好的payload
	ch := parser.ParseStream(conn)
	for payload := range ch {
		if payload.Err != nil {
			if payload.Err == io.EOF ||
				payload.Err == io.ErrUnexpectedEOF ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {
				// 断开连接，终止处理
				h.closeClient(client)
				slog.Info("connection closed: " + client.RemoteAddr())
				return
			}

			errReply := protocol.MakeErrReply(payload.Err.Error())
			// 往客户端写一个错误回复
			_, err := client.Write(errReply.ToBytes())
			if err != nil {
				h.closeClient(client)
				slog.Info("connection closed: " + client.RemoteAddr())
				return
			}
			continue
		}
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
		result := h.dbEngine.Exec(client, r.Args)
		if result != nil {
			_, _ = client.Write(result.ToBytes())
		} else {
			_, _ = client.Write([]byte("-ERR unknown\r\n"))
		}
	}
}
