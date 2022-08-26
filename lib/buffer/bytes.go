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

func Write8(buf io.Writer, val uint8) error {
	tSize := 1
	data := make([]byte, tSize)
	_, err := buf.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func LittleEndianWrite16(buf io.Writer, val uint16) error {
	tSize := 2
	data := make([]byte, tSize)
	binary.LittleEndian.PutUint16(data, val)
	_, err := buf.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func LittleEndianWrite32(buf io.Writer, val uint32) error {
	tSize := 4
	data := make([]byte, tSize)
	binary.LittleEndian.PutUint32(data, val)
	_, err := buf.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func LittleEndianWrite64(buf io.Writer, val uint64) error {
	tSize := 8
	data := make([]byte, tSize)
	binary.LittleEndian.PutUint64(data, val)
	_, err := buf.Write(data)
	if err != nil {
		return err
	}
	return nil
}
