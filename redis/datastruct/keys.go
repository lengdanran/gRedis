// Package datastruct lengdanran 2024/4/2 18:46
package datastruct

import (
	"github.com/lengdanran/gredis/interface/redis"
	"github.com/lengdanran/gredis/lib/wildcard"
	"github.com/lengdanran/gredis/redis/dbengine"
	"github.com/lengdanran/gredis/redis/protocol"
)

func init() {
	dbengine.RegisterExecutor("keys", ExeKeys)
	dbengine.RegisterExecutor("exists", ExeExists)
}

func ExeKeys(eg *dbengine.RedisEngine, args [][]byte) redis.Reply {
	pattern, err := wildcard.CompilePattern(string(args[0]))
	if err != nil {
		return protocol.MakeErrReply("ERR illegal wildcard")
	}
	result := make([][]byte, 0)
	keys := eg.Data.Keys()
	for _, k := range keys {
		if !pattern.IsMatch(k) {
			continue
		}
		if eg.IsExpired(k) {
			continue
		}
		result = append(result, []byte(k))
	}
	return protocol.MakeMultiBulkReply(result)
}

func ExeExists(eg *dbengine.RedisEngine, args [][]byte) redis.Reply {
	contains := eg.Data.Contains(string(args[0]))
	if contains {
		return protocol.MakeBulkReply([]byte("true"))
	}
	return protocol.MakeBulkReply([]byte("false"))
}
