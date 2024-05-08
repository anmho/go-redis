package redis

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type RespDataType string

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

func NewString(s string) Value {
	return Value{typ: "string", str: s}
}

func NewInteger(num int) Value {
	return Value{typ: "integer", num: num}
}
func NewBulkString(s string) Value {
	return Value{typ: "bulk", bulk: s}
}

func NewArray(a []Value) Value {
	return Value{typ: "array", array: a}
}

func (v Value) Marshal() []byte {
	switch v.typ {

	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "null":
		return v.marshalNull()
	case "int":
		return v.marshalInteger()
	case "error":
		return v.marshalError()
	default:
		return []byte{}
	}
}

func (v Value) marshalArray() []byte {
	var bytes []byte
	numElements := len(v.array)
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(numElements)...)
	bytes = append(bytes, '\r', '\n')
	for _, child := range v.array {
		//fmt.Println((child.Marshal()))
		bytes = append(bytes, child.Marshal()...)
	}
	return bytes
}

func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}
func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshalInteger() []byte {
	var bytes []byte
	bytes = append(bytes, INTEGER)
	bytes = append(bytes, strconv.Itoa(v.num)...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}

func (v Value) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

// readLine translates the bytes from in the reader
func (r *Resp) readLine() (line []byte, n int, err error) {
	// do while to capture all bytes read
	// since its not null terminated
	for {

		// read the byte
		b, err := r.reader.ReadByte()

		if err != nil {
			return nil, 0, err
		}
		// n is the num bytes read
		n += 1
		line = append(line, b) // build the line

		// stop if the line length is only 2
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}

	}
	// n is the size of the message without the crlf
	return line[:len(line)-2], n, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte() // read the first byte to identify the input
	if err != nil {
		return Value{}, fmt.Errorf("reading type byte: %w", err)
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, nil
	}
}

func (r *Resp) readArray() (Value, error) {
	// at this point we just figured out that the type is an array bc of the symbol
	var v Value
	v.typ = "array"
	length, _, err := r.readInteger() // how many elems in the array
	if err != nil {
		return v, err
	}
	log.Println("length:", length)
	v.array = make([]Value, 0)
	// do the same for every element in the array
	//for i := 0; i < length; i++ {
	for range length {
		val, err := r.Read()
		if err != nil {
			return v, err
		}
		v.array = append(v.array, val)

	}
	return v, nil
}

func (r *Resp) readBulk() (Value, error) {
	/* original string before reading symbol
		$5\r\nhello\r\n
	eg
		$5
		hello
	*/

	v := Value{}
	v.typ = "bulk"

	length, _, err := r.readInteger()
	if err != nil {
		return Value{}, fmt.Errorf("reading length of bulk: %w", err)
	}

	bulk := make([]byte, length)
	_, err = r.reader.Read(bulk)
	if err != nil {
		return Value{}, fmt.Errorf("reading bulk of %d bytes: %w", length, err)
	} // fills the buffer?
	v.bulk = string(bulk)
	_, _, err = r.readLine()
	if err != nil {
		return Value{}, fmt.Errorf("reading remaining CRLF of bulk: %w", err)
	}

	return v, nil

}
