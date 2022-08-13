package lib

import (
	"encoding/binary"
	"io"
)

func Read8(buf io.ByteReader) (uint8, error) {
	b, err := buf.ReadByte()
	if err != nil {
		return 0, nil
	}
	return b, nil
}

func Read16(buf io.Reader) (uint16, error) {
	b := make([]byte, 2)
	_, err := buf.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(b), nil
}

func Read32(buf io.Reader) (uint32, error) {
	b := make([]byte, 4)
	_, err := buf.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(b), nil
}

func Read64(buf io.Reader) (uint64, error) {
	b := make([]byte, 8)
	_, err := buf.Read(b)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(b), nil
}

func Write8(buf io.ByteWriter, val uint8) error {
	return buf.WriteByte(val)
}

func Write16(buf io.Writer, val uint16) error {
	tSize := 2
	data := make([]byte, tSize)
	binary.LittleEndian.PutUint16(data, val)
	_, err := buf.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func Write32(buf io.Writer, val uint32) error {
	tSize := 4
	data := make([]byte, tSize)
	binary.LittleEndian.PutUint32(data, val)
	_, err := buf.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func Write64(buf io.Writer, val uint64) error {
	tSize := 8
	data := make([]byte, tSize)
	binary.LittleEndian.PutUint64(data, val)
	_, err := buf.Write(data)
	if err != nil {
		return err
	}
	return nil
}
