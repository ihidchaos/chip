package buffer

import (
	"encoding/binary"
	"io"
)

func Write8(buf io.Writer, val uint8) error {
	data := []byte{val}
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
