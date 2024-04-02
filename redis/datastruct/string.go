// Package datastruct lengdanran 2024/4/1 10:26
package datastruct

import (
	"github.com/lengdanran/gredis/interface/redis"
	"github.com/lengdanran/gredis/redis/dbengine"
	"github.com/lengdanran/gredis/redis/protocol"
	"log"
)

func init() {
	// register string executor options into ExecutorMap
	dbengine.RegisterExecutor("get", ExeGet)
	dbengine.RegisterExecutor("set", ExeSet)
	dbengine.RegisterExecutor("del", ExeDel)
}

func getAsString(eg *dbengine.RedisEngine, key string) ([]byte, protocol.ErrorReply) {
	entity, ok := eg.Get(key)
	if !ok {
		return nil, nil
	}
	bytes, ok := entity.Data.([]byte)
	if !ok {
		return nil, &protocol.WrongTypeErrReply{}
	}
	return bytes, nil
}

func ExeGet(eg *dbengine.RedisEngine, args [][]byte) redis.Reply {
	key := string(args[0])
	bytes, err := getAsString(eg, key)
	if err != nil {
		log.Print(err.Error())
		return protocol.MakeErrReply(err.Error())
	}
	if bytes == nil {
		return &protocol.NullBulkReply{}
	}
	return protocol.MakeBulkReply(bytes)
}

func ExeSet(eg *dbengine.RedisEngine, args [][]byte) redis.Reply {
	key := string(args[0])
	value := args[1]
	eg.PutEntity(key, &dbengine.DataEntity{Data: value})
	return &protocol.OkReply{}
}

func ExeDel(eg *dbengine.RedisEngine, args [][]byte) redis.Reply {
	key := string(args[0])
	eg.Remove(key)
	return &protocol.OkReply{}
}
