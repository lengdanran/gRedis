// Package dbengine lengdanran 2024/3/28 10:13
package dbengine

import (
	"fmt"
	"github.com/lengdanran/gredis/interface/redis"
	"github.com/lengdanran/gredis/redis/connection"
	"github.com/lengdanran/gredis/redis/datastruct/dict"
	"github.com/lengdanran/gredis/redis/protocol"
	"log/slog"
	"runtime/debug"
	"strings"
)

type RedisEngine struct {
	DBEngine
	// key -> DataEntity
	data *dict.ConcurrentDict
	// key -> expireTime (time.Time)
	ttlMap *dict.ConcurrentDict
	// key -> version(uint32)
	versionMap *dict.ConcurrentDict
	// 回调钩子
	insertCallback KeyEventCallback
	deleteCallback KeyEventCallback
}

func NewRedisEngine() *RedisEngine {
	engine := &RedisEngine{}
	return engine
}

func (engine *RedisEngine) SetInsertCallback(callback KeyEventCallback) {
	engine.insertCallback = callback
}

func (engine *RedisEngine) SetDeleteCallback(callback KeyEventCallback) {
	engine.deleteCallback = callback
}

/* ---- Lock Function ----- */

// RWLocks lock keys for writing and reading
func (engine *RedisEngine) RWLocks(writeKeys []string, readKeys []string) {
	engine.data.RWLocks(writeKeys, readKeys)
}

// RWUnLocks unlock keys for writing and reading
func (engine *RedisEngine) RWUnLocks(writeKeys []string, readKeys []string) {
	engine.data.RWUnLocks(writeKeys, readKeys)
}

func (engine *RedisEngine) Exec(cmdLine [][]byte) (result redis.Reply) {
	defer func() {
		if err := recover(); err != nil {
			slog.Warn(fmt.Sprintf("error occurs: %v\n%s", err, string(debug.Stack())))
			result = &protocol.UnknownErrReply{}
		}
	}()
	exeName := strings.ToLower(string(cmdLine[0]))
	etr, ok := ExecutorMap[exeName]
	if !ok {
		return protocol.MakeErrReply("ERR unknown command '" + exeName + "'")
	}
	write, read := etr.LockKeysF(cmdLine[1:])
	engine.RWLocks(write, read)
	defer engine.RWUnLocks(write, read)
	result = etr.ExecF(engine, cmdLine[1:])
	return result
}

func (engine *RedisEngine) AfterClientClose(c *connection.Connection) {
	// do nothing
}

func (engine *RedisEngine) Close() {
	// do nothing
}

func (engine *RedisEngine) GetEntity(key string) (*DataEntity, bool) {
	return nil, false
}
