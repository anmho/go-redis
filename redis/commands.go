package redis

import "errors"

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
}

//var Handlers = map[string]func([]Value) Value

func ping(args []Value) Value {
	return NewString("PONG")
}

func Handle(v Value) error {
	if v.typ != "array" {
		return errors.New("want array")
	}

	return nil

}
