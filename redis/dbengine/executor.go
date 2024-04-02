// Package dbengine lengdanran 2024/3/28 11:21
package dbengine

import (
	"log"
	"strings"
)

var ExecutorMap = make(map[string]*Executor)

type Executor struct {
	Name  string
	ExecF ExecFunc
}

// RegisterExecutor registers a normal command, which only read or modify a limited number of keys
func RegisterExecutor(name string, f ExecFunc) {
	log.Print("RegisterExecutor >> ", name)
	name = strings.ToLower(name)
	executor := &Executor{
		Name:  name,
		ExecF: f,
	}
	ExecutorMap[name] = executor
}
