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
