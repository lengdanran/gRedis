// Package dbengine lengdanran 2024/3/27 19:40
package dbengine

import (
	"github.com/lengdanran/gredis/interface/redis"
)

type DBEngine interface {
	Exec(cmdLine [][]byte) redis.Reply
}

type ExecFunc func(eg *RedisEngine, args [][]byte) redis.Reply

// DataEntity 代表数据的存储方式，键值存储，string, list, hash, set and so on
type DataEntity struct {
	Data interface{}
}
