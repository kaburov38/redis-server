package main

import (
	"fmt"
	"sync"
)

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
}
var Sets = map[string]string{}
var SetsMu = sync.RWMutex{}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}

	return Value{typ: "string", str: args[0].bulk}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "Wrong number of arguments for SET"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SetsMu.Lock()
	Sets[key] = value
	SetsMu.Unlock()
	fmt.Println(Sets)

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "Wrong number of arguments for GET"}
	}

	key := args[0].bulk
	fmt.Println(key)
	SetsMu.RLock()
	value, ok := Sets[key]
	SetsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}
