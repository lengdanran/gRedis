// Package dbengine lengdanran 2024/3/27 19:40
package dbengine

import (
	"github.com/lengdanran/gredis/interface/redis"
	"github.com/lengdanran/gredis/redis/connection"
	"github.com/lengdanran/gredis/redis/protocol"
)

type DBEngine interface {
	Exec(client *connection.Connection, cmdLine [][]byte) redis.Reply
	AfterClientClose(c *connection.Connection)
	Close()
	// LoadRDB(dec *core.Decoder) error
}

type SimpleEngine struct {
}

func (engine *SimpleEngine) Exec(client *connection.Connection, cmdLine [][]byte) redis.Reply {
	return protocol.MakeBulkReply([]byte("msg from SimpleEngine"))
}

func (engine *SimpleEngine) AfterClientClose(c *connection.Connection) {
	// do nothing
}

func (engine *SimpleEngine) Close() {
	// do nothing
}
