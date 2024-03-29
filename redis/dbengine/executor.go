// Package executor lengdanran 2024/3/28 11:21
package dbengine

import (
	"strings"
)

var ExecutorMap = make(map[string]*Executor)

type Executor struct {
	Name      string
	ExecF     ExecFunc
	LockKeysF GetLockKeys
}

// RegisterExecutor registers a normal command, which only read or modify a limited number of keys
func RegisterExecutor(name string, f ExecFunc, getLockKeysF GetLockKeys) *Executor {
	name = strings.ToLower(name)
	executor := &Executor{
		Name:      name,
		ExecF:     f,
		LockKeysF: getLockKeysF,
	}
	ExecutorMap[name] = executor
	return executor
}
