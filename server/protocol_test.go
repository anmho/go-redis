package server

import (
	"bytes"
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
	assert.Equal(4, n)
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

//func TestResp_readArray(t *testing.T) {
//	//assert := assert.New(t)
//
//}
