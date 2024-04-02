// Package datastruct lengdanran 2024/4/1 10:26
package datastruct

import (
	"github.com/lengdanran/gredis/interface/redis"
	"github.com/lengdanran/gredis/redis/dbengine"
	"github.com/lengdanran/gredis/redis/protocol"
)

func init() {
	// register string executor options into ExecutorMap
	dbengine.RegisterExecutor("get", ExeGet)
	dbengine.RegisterExecutor("set", ExeSet)
	dbengine.RegisterExecutor("del", ExeDel)
	dbengine.RegisterExecutor("getset", ExeGetSet)
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
	bytes, errReply := getAsString(eg, key)
	if errReply != nil {
		return errReply
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

func ExeGetSet(eg *dbengine.RedisEngine, args [][]byte) redis.Reply {
	if len(args) < 2 {
		return protocol.MakeArgNumErrReply("getset cmd needs 2 arguments!")
	}
	key := string(args[0])
	val := string(args[1])
	oldVal := eg.Data.Get(key).(string)
	eg.PutEntity(key, &dbengine.DataEntity{Data: val})
	return protocol.MakeBulkReply([]byte(oldVal))
}
