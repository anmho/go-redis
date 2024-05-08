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
		input Value
		want  []byte
	}{
		{
			Value{"string", "OK", 0, "", nil},
			[]byte("+OK\r\n"),
		},
		{
			input: Value{"error", "Error message", 0, "", nil},
			want:  []byte("-Error message\r\n"),
		},
		{
			input: Value{"int", "", 10, "", nil},
			want:  []byte(":10\r\n"),
		},
		{
			input: Value{"bulk", "", 0, "bulk-string", nil},
			want:  []byte("$11\r\nbulk-string\r\n"),
		},
		{
			input: Value{"bulk", "", 0, "hello", nil},
			want:  []byte("$5\r\nhello\r\n"),
		},
		{
			input: Value{"bulk", "", 0, "", nil},
			want:  []byte("$0\r\n\r\n"),
		},
	}

	for i, test := range tests {
		log.Println("Case", i)
		assert.Equal(string(test.want), string(test.input.Marshal()))
	}
}

//func TestResp_readArray(t *testing.T) {
//	//assert := assert.New(t)
//
//}
