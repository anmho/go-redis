package redis

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
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
	t.Parallel()
	assert := assert.New(t)

	v := Value{
		typ:   StringType,
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
	t.Parallel()
	tests := []struct {
		input       Value
		want        []byte
		description string
	}{
		{
			Value{typ: StringType, str: "OK"},
			[]byte("+OK\r\n"),
			"OK",
		},
		{
			Value{typ: ErrorType, str: "Error message"},
			[]byte("-Error message\r\n"),
			"error",
		},
		{
			Value{typ: IntType, num: 10},
			[]byte(":10\r\n"),
			"number",
		},
		{
			Value{typ: BulkType, bulk: "bulk-string"},
			[]byte("$11\r\nbulk-string\r\n"),
			"bulk string",
		},
		{
			Value{typ: BulkType, bulk: "hello"},
			[]byte("$5\r\nhello\r\n"),
			"hello",
		},
		{
			Value{typ: BulkType, bulk: ""},
			[]byte("$0\r\n\r\n"),
			"empty string",
		},
		{
			Value{typ: ArrayType, array: []Value{
				{typ: StringType, str: "A"},
				{typ: StringType, str: "B"},
				{typ: StringType, str: "C"},
			}},
			[]byte("*3\r\n+A\r\n+B\r\n+C\r\n"),
			"array",
		},
	}

	for i, test := range tests {
		name := fmt.Sprintf("case %d: %s", i, test.description)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, string(test.want), string(test.input.Marshal()))
		})

	}
}
