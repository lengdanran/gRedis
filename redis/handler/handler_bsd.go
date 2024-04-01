// Package handler lengdanran 2024/3/27 18:43
//go:build darwin

package handler

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/lengdanran/gredis/epoll"
	"github.com/lengdanran/gredis/redis/dbengine"
	"github.com/lengdanran/gredis/redis/parser"
	"github.com/lengdanran/gredis/redis/protocol"
	"io"
	"log/slog"
	"strings"
)

func RedisBsdPollEventHandleFunc(s *epoll.BsdSocket, engine dbengine.DBEngine) {
	reader := bufio.NewReader(s)
	ch := make(chan *parser.Payload)
	go parser.ParseRedisRequestStream(reader, ch)
	for payload := range ch {
		if payload.Err != nil {
			if payload.Err == io.EOF ||
				errors.Is(payload.Err, io.ErrUnexpectedEOF) ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {
				// 断开连接，终止处理
				slog.Info("connection closed!")
				return
			}
			errReply := protocol.MakeErrReply(payload.Err.Error())
			_, err := s.Write(errReply.ToBytes())
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
				_, _ = s.Write(result.ToBytes())
			} else {
				_, _ = s.Write([]byte("-ERR unknown\r\n"))
			}
		}
	}
}
