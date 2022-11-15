package spake2

import (
	"bytes"
	"encoding/binary"
)

func concat(bytesArray ...[]byte) []byte {
	var result []byte
	for _, b := range bytesArray {
		if len(b) > 0 {
			bytesLen := make([]byte, 8)
			binary.LittleEndian.PutUint64(bytesLen, uint64(len(b)))
			result = append(result, bytesLen...)
			result = append(result, b...)
		}
	}
	return result
}

func padScalarBytes(scBytes []byte, padLen int) []byte {
	if len(scBytes) > padLen {
		return scBytes
	}
	return append(bytes.Repeat([]byte{0}, padLen-len(scBytes)), scBytes...)
}

func appendLenAndContent(b *bytes.Buffer, input []byte) {
	binary.Write(b, binary.LittleEndian, uint64(len(input)))
	b.Write(input)
}
