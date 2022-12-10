package secure_channel

import (
	"bytes"
	"encoding/binary"
	"github.com/galenliu/chip/protocols"
	"io"
)

type StatusReport struct {
	ProtocolCode uint16
	ProtocolId   protocols.Id
	GeneralCode  generalStatusCode
	ProtocolData []byte
}

func (s *StatusReport) Decode(buf *bytes.Buffer) (err error) {

	tempGeneralCode := binary.LittleEndian.Uint16(buf.Next(2))

	tProtocolId := binary.LittleEndian.Uint32(buf.Next(4))

	s.ProtocolCode = binary.LittleEndian.Uint16(buf.Next(2))

	s.ProtocolId = protocols.FromFullyQualifiedSpecForm(tProtocolId)

	s.GeneralCode = generalStatusCode(tempGeneralCode)

	if buf.Len() > 0 {
		s.ProtocolData = buf.Bytes()
	} else {
		s.ProtocolData = nil
	}
	return err
}

func (s *StatusReport) Encode(buf io.Writer) (err error) {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, uint16(s.GeneralCode))
	if _, err = buf.Write(data); err != nil {
		return err
	}
	data = make([]byte, 4)
	binary.LittleEndian.PutUint32(data, s.ProtocolId.ToFullyQualifiedSpecForm())
	if _, err = buf.Write(data); err != nil {
		return err
	}

	data = make([]byte, 2)
	binary.LittleEndian.PutUint16(data, s.ProtocolCode)
	if _, err = buf.Write(data); err != nil {
		return err
	}

	if s.ProtocolData != nil {
		if _, err := buf.Write(s.ProtocolData); err != nil {
			return err
		}
	}
	return nil
}
