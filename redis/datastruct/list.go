// Package datastruct lengdanran 2024/4/7 16:47
package datastruct

import (
	"github.com/lengdanran/gredis/interface/redis"
	L "github.com/lengdanran/gredis/lib/list"
	"github.com/lengdanran/gredis/lib/utils"
	"github.com/lengdanran/gredis/redis/dbengine"
	"github.com/lengdanran/gredis/redis/protocol"
	"strconv"
)

func init() {
	// register list executor options into ExecutorMap
	dbengine.RegisterExecutor("lpush", ExeLPush)
	dbengine.RegisterExecutor("rpush", ExeRPush)
	dbengine.RegisterExecutor("lpop", ExeLPop)
	dbengine.RegisterExecutor("rpop", ExeRPop)
	dbengine.RegisterExecutor("llen", ExeLLen)
	dbengine.RegisterExecutor("lrange", ExeLRange)
	dbengine.RegisterExecutor("lrem", ExeLRem)
}

func getAsList(eg *dbengine.RedisEngine, key string) (L.List, protocol.ErrorReply) {
	entity, _ := eg.Get(key)
	if entity == nil {
		lst := L.NewQuickList()
		eg.PutEntity(key, &dbengine.DataEntity{
			Data: lst,
		})
		return lst, nil
	} else {
		lst, ok := entity.Data.(L.List)
		if !ok {
			return nil, &protocol.WrongTypeErrReply{}
		}
		return lst, nil
	}
}

func ExeLPush(eg *dbengine.RedisEngine, args [][]byte) redis.Reply {
	key := string(args[0])
	values := args[1:]
	lst, errReply := getAsList(eg, key)
	if errReply != nil {
		return errReply
	}
	for _, value := range values {
		lst.Insert(0, value)
	}
	return protocol.MakeIntReply(int64(lst.Len()))
}

func ExeRPush(eg *dbengine.RedisEngine, args [][]byte) redis.Reply {
	key := string(args[0])
	values := args[1:]
	lst, errReply := getAsList(eg, key)
	if errReply != nil {
		return errReply
	}
	for _, value := range values {
		lst.Add(value)
	}
	return protocol.MakeIntReply(int64(lst.Len()))
}

func ExeLPop(eg *dbengine.RedisEngine, args [][]byte) redis.Reply {
	key := string(args[0])
	lst, errReply := getAsList(eg, key)
	if errReply != nil {
		return errReply
	}
	bytes := lst.Remove(0).([]byte)
	if lst.Len() == 0 {
		eg.Data.Del(key)
	}
	return protocol.MakeBulkReply(bytes)
}

func ExeRPop(eg *dbengine.RedisEngine, args [][]byte) redis.Reply {
	key := string(args[0])
	lst, errReply := getAsList(eg, key)
	if errReply != nil {
		return errReply
	}
	if lst.Len() == 0 {
		eg.Data.Del(key)
		return protocol.MakeNullBulkReply()
	}
	last := lst.RemoveLast()
	var bytes []byte
	if last == nil {
		bytes = []byte("empty")
	} else {
		bytes = last.([]byte)
	}
	return protocol.MakeBulkReply(bytes)
}

func ExeLLen(eg *dbengine.RedisEngine, args [][]byte) redis.Reply {
	key := string(args[0])
	lst, errReply := getAsList(eg, key)
	if errReply != nil {
		return errReply
	}
	return protocol.MakeIntReply(int64(lst.Len()))
}

// ExeLRange gets elements of list in given range
func ExeLRange(eg *dbengine.RedisEngine, args [][]byte) redis.Reply {
	key := string(args[0])
	start64, err := strconv.ParseInt(string(args[1]), 10, 64)
	if err != nil {
		return protocol.MakeErrReply("ERR value is not an integer or out of range")
	}
	start := int(start64)
	stop64, err := strconv.ParseInt(string(args[2]), 10, 64)
	if err != nil {
		return protocol.MakeErrReply("ERR value is not an integer or out of range")
	}
	stop := int(stop64)

	// get data
	lst, errReply := getAsList(eg, key)
	if errReply != nil {
		return errReply
	}
	if lst == nil {
		return &protocol.EmptyMultiBulkReply{}
	}

	// compute index
	size := lst.Len() // assert: size > 0
	if start < -1*size {
		start = 0
	} else if start < 0 {
		start = size + start
	} else if start >= size {
		return &protocol.EmptyMultiBulkReply{}
	}
	if stop < -1*size {
		stop = 0
	} else if stop < 0 {
		stop = size + stop + 1
	} else if stop < size {
		stop = stop + 1
	} else {
		stop = size
	}
	if stop < start {
		stop = start
	}

	// assert: start in [0, size - 1], stop in [start, size]
	slice := lst.Range(start, stop)
	result := make([][]byte, len(slice))
	for i, raw := range slice {
		bytes, _ := raw.([]byte)
		result[i] = bytes
	}
	return protocol.MakeMultiBulkReply(result)
}

// ExeLRem removes element of list at specified index
func ExeLRem(eg *dbengine.RedisEngine, args [][]byte) redis.Reply {
	// parse args
	key := string(args[0])
	count64, err := strconv.ParseInt(string(args[1]), 10, 64)
	if err != nil {
		return protocol.MakeErrReply("ERR value is not an integer or out of range")
	}
	count := int(count64)
	value := args[2]

	// get data entity
	list, errReply := getAsList(eg, key)
	if errReply != nil {
		return errReply
	}
	if list == nil {
		return protocol.MakeIntReply(0)
	}

	var removed int
	if count == 0 {
		removed = list.RemoveAllByVal(func(a interface{}) bool {
			return utils.Equals(a, value)
		})
	} else if count > 0 {
		removed = list.RemoveByVal(func(a interface{}) bool {
			return utils.Equals(a, value)
		}, count)
	} else {
		removed = list.ReverseRemoveByVal(func(a interface{}) bool {
			return utils.Equals(a, value)
		}, -count)
	}

	if list.Len() == 0 {
		eg.Data.Del(key)
	}

	return protocol.MakeIntReply(int64(removed))
}
