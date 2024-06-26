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

type RespDataType = int

const (
	ArrayType RespDataType = iota + 1
	BulkType
	StringType
	NullType
	IntType
	ErrorType
)

type Value struct {
	typ   RespDataType
	str   string
	num   int
	bulk  string
	array []Value
}

func NewString(s string) Value {
	return Value{typ: StringType, str: s}
}

func NewInteger(num int) Value {
	return Value{typ: IntType, num: num}
}
func NewBulkString(s string) Value {
	return Value{typ: BulkType, bulk: s}
}

func NewArray(a []Value) Value {
	return Value{typ: ArrayType, array: a}
}

func NewError(format string, values ...any) Value {
	return Value{typ: ErrorType, str: fmt.Sprintf(format, values...)}
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
	v.typ = ArrayType
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
	v.typ = BulkType

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
