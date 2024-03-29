// Package dbengine lengdanran 2024/3/27 19:40
package dbengine

import (
	"github.com/lengdanran/gredis/interface/redis"
	"github.com/lengdanran/gredis/redis/connection"
	"github.com/lengdanran/gredis/redis/protocol"
)

type DBEngine interface {
	Exec(cmdLine [][]byte) redis.Reply
	AfterClientClose(c *connection.Connection)
	Close()
	// LoadRDB(dec *core.Decoder) error
}

type ExecFunc func(eg *RedisEngine, args [][]byte) redis.Reply
type GetLockKeys func(args [][]byte) ([]string, []string)

// DataEntity 代表数据的存储方式，键值存储，string, list, hash, set and so on
type DataEntity struct {
	Data interface{}
}

type CmdLine = [][]byte

// KeyEventCallback 操作完成的回调
// may be called concurrently
type KeyEventCallback func(dbIndex int, key string, entity *DataEntity)

type SimpleEngine struct {
}

func (engine *SimpleEngine) Exec(cmdLine [][]byte) redis.Reply {
	return protocol.MakeBulkReply([]byte("msg from SimpleEngine"))
}

func (engine *SimpleEngine) AfterClientClose(c *connection.Connection) {
	// do nothing
}

func (engine *SimpleEngine) Close() {
	// do nothing
}
