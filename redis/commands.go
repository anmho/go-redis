package redis

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
)

var (
	ErrBadRequest   = errors.New("invalid request")
	ErrRespProtocol = errors.New("invalid resp message")
	ErrInvalidCmd   = errors.New("invalid command")
)

var handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func Handle(req Value) (Value, error) {
	if req.typ != ArrayType {
		return Value{}, fmt.Errorf("expected array type: %w", ErrBadRequest)
	}
	if len(req.array) == 0 {
		return Value{}, fmt.Errorf("expected non-empty array: %w", ErrBadRequest)
	}
	if req.array[0].typ != BulkType {
		return Value{}, fmt.Errorf("expected first element bulk type: %w", ErrBadRequest)
	}
	command := strings.ToUpper(req.array[0].bulk)
	log.Println("executing", command)
	handler, ok := handlers[command]
	if !ok {
		return Value{}, fmt.Errorf("unknown command '%s': %w", command, ErrBadRequest)
	}
	args := req.array[1:]
	result := handler(args)
	return result, nil
}

func ping(args []Value) Value {
	return NewString("PONG")
}

func set(args []Value) Value {
	if len(args) != 2 {
		return NewError(fmt.Sprintf("expected 2 arguments got %d", len(args)))
	}

	var key = args[0].bulk
	var value = args[1].bulk

	// Lock for everyone else
	SETsMu.Lock()
	defer SETsMu.Unlock()
	SETs[key] = value

	return NewString("OK")
}

func get(args []Value) Value {
	if len(args) != 1 {
		return NewError("expected 1 argument, got %d", len(args))
	}
	// Lock for writing
	// Allow writing
	SETsMu.RLock()
	defer SETsMu.RUnlock()
	var key = args[0].bulk
	var value = SETs[key]
	return NewString(value)
}

func hset(args []Value) Value {
	if len(args) != 3 {
		return NewError("expected 3 args, got", len(args))
	}

	var key = args[0].bulk
	var field = args[1].bulk
	var value = args[2].bulk
	HSETsMu.Lock()
	defer HSETsMu.Unlock()
	if _, ok := HSETs[key]; !ok {
		HSETs[key] = map[string]string{}
	}
	HSETs[key][field] = value
	return NewString("OK")
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return NewError("expected 3 arg, got", len(args))
	}

	var key = args[0].bulk
	var field = args[1].bulk

	HSETsMu.RLock()
	defer HSETsMu.RUnlock()
	var value = HSETs[key][field]

	return NewString(value)
}

func hgetall(args []Value) Value {
	log.Println("getting all values associated with key")

	return Value{}
}
