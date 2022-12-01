package tlv

import "reflect"

type errorContext struct {
	Struct     reflect.Type
	FieldStack []string
}

// decodeState represents the state while decoding a JSON value.
type decodeState struct {
	data                  []byte
	off                   int // next read offset in data
	opcode                int // last read result
	scan                  scanner
	errorContext          *errorContext
	savedError            error
	useNumber             bool
	disallowUnknownFields bool
}

func (d *decodeState) init(data []byte) *decodeState {
	d.data = data
	d.off = 0
	d.savedError = nil
	if d.errorContext != nil {
		d.errorContext.Struct = nil
		// Reuse the allocated space for the FieldStack slice.
		d.errorContext.FieldStack = d.errorContext.FieldStack[:0]
	}
	return d
}
