package redis

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestResp_readLine(t *testing.T) {
	assert := assert.New(t)

	reader := bytes.NewReader([]byte("5\r\n"))
	r := NewResp(reader)
	line, n, err := r.readLine()
	assert.NoError(err)

	assert.Equal("5", string(line))
	assert.Equal(3, n)
}

func TestResp_readInteger(t *testing.T) {
	assert := assert.New(t)
	reader := bytes.NewReader([]byte("0\r\n")) // expects the first byte to already be read
	r := NewResp(reader)

	x, n, err := r.readInteger()
	assert.NoError(err)
	assert.Equal(0, x)
	assert.Equal(3, n)
}

func TestResp_marshalString(t *testing.T) {
	assert := assert.New(t)

	v := Value{
		typ:   "string",
		str:   "hello",
		num:   0,
		bulk:  "",
		array: nil,
	}
	data := v.Marshal()
	assert.Equal("+hello\r\n", string(data))

}

func Test_marshalArray(t *testing.T) {

}

func TestResp_Marshal(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		input       Value
		want        []byte
		description string
	}{
		{
			Value{typ: "string", str: "OK"},
			[]byte("+OK\r\n"),
			"OK",
		},
		{
			Value{typ: "error", str: "Error message"},
			[]byte("-Error message\r\n"),
			"error",
		},
		{
			Value{typ: "int", num: 10},
			[]byte(":10\r\n"),
			"number",
		},
		{
			Value{typ: "bulk", bulk: "bulk-string"},
			[]byte("$11\r\nbulk-string\r\n"),
			"bulk string",
		},
		{
			Value{typ: "bulk", bulk: "hello"},
			[]byte("$5\r\nhello\r\n"),
			"hello",
		},
		{
			Value{typ: "bulk", bulk: ""},
			[]byte("$0\r\n\r\n"),
			"empty string",
		},
		{
			Value{typ: "array", array: []Value{
				{typ: "string", str: "A"},
				{typ: "string", str: "B"},
				{typ: "string", str: "C"},
			}},
			[]byte("*3\r\n+A\r\n+B\r\n+C\r\n"),
			"array",
		},
	}

	for i, test := range tests {
		log.Println("Case", i, test.description)
		assert.Equal(string(test.want), string(test.input.Marshal()))
	}
}
