package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

type Resp struct {
	reader *bufio.Reader
}

type Writer struct {
	writer io.Writer
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()

		if err != nil {
			if err == io.EOF {
				return line, n, nil
			}

			return nil, 0, err
		}

		if b == '\r' {
			r.reader.ReadLine()
			break
		}

		n += 1
		line = append(line, b)
	}

	return line, n, err
}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()

	if err != nil {
		return 0, 0, err
	}

	num, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return int(num), n, nil
}

func (r *Resp) readArray() (v Value, err error) {
	v.typ = "array"
	length, _, err := r.readInteger()
	if err != nil {
		fmt.Println(err)
		return v, err
	}
	v.array = make([]Value, length)
	for i := 0; i < length; i++ {
		_value, err := r.Read()
		if err != nil {
			fmt.Println(err)
			return v, err
		}
		v.array[i] = _value
	}
	return
}

func (r *Resp) readBulk() (v Value, err error) {
	v.typ = "bulk"

	_length, _, err := r.readInteger()
	if err != nil {
		fmt.Println(err)
		return v, err
	}

	result := make([]byte, _length)

	r.reader.Read(result)

	v.bulk = string(result)

	//clean up the \r\n
	r.readLine()

	return
}

func (r *Resp) Read() (v Value, err error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		fmt.Println(err)
		return
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Println("Unknow Type: ", string(_type))
		return Value{}, nil
	}
}

func (value Value) Marshal() []byte {
	switch value.typ {
	case "array":
		return value.marshalArray()
	case "bulk":
		return value.marshalBulk()
	case "string":
		return value.marshalString()
	case "null":
		return value.marshalNull()
	case "error":
		return value.marshalError()
	default:
		return []byte{}
	}
}

func (value Value) marshalArray() []byte {
	len := len(value.array)

	var bytes []byte

	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')

	for i := range len {
		bytes = append(bytes, value.array[i].Marshal()...)
	}
	return bytes
}

func (value Value) marshalBulk() []byte {
	len := len(value.bulk)

	var bytes []byte

	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, value.bulk...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (value Value) marshalString() []byte {
	var bytes []byte

	bytes = append(bytes, STRING)
	bytes = append(bytes, value.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (value Value) marshalError() []byte {
	var bytes []byte

	bytes = append(bytes, ERROR)
	bytes = append(bytes, value.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (value Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (w Writer) Write(value Value) error {
	bytes := value.Marshal()

	_, err := w.writer.Write(bytes)
	return err
}
