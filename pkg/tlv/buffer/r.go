package buffer

import (
	"encoding/binary"
	"io"
)

func Read8(buf io.Reader) (uint8, error) {
	data := make([]byte, 1)
	_, err := buf.Read(data)
	if err != nil {
		return 0, nil
	}
	return data[0], nil
}

func LittleEndianRead16(buf io.Reader) (uint16, error) {
	b := make([]byte, 2)
	_, err := buf.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(b), nil
}

func LittleEndianRead32(buf io.Reader) (uint32, error) {
	b := make([]byte, 4)
	_, err := buf.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(b), nil
}

func LittleEndianRead64(buf io.Reader) (uint64, error) {
	b := make([]byte, 8)
	_, err := buf.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(b), nil
}
