package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Header struct {
	code  byte
	size  byte
	hType byte
}

func ByteToHeader(buffer *[]byte, header *Header) error {
	if header == nil {
		return errors.New("header is nil")
	}

	msgBuf := bytes.NewBuffer(*buffer)

	binary.Read(msgBuf, binary.LittleEndian, &header.code)
	if header.code != 0x89 {
		return errors.New("not equal code")
	}

	if msgBuf.Len() < 8 {
		return errors.New("length is less copy")
	}
	binary.Read(msgBuf, binary.LittleEndian, &header.size)
	binary.Read(msgBuf, binary.LittleEndian, &header.hType)
	return nil
}

func HeaderToByte(header *Header, buffer *bytes.Buffer) int {
	prev := buffer.Len()

	binary.Write(buffer, binary.LittleEndian, uint64(0x89))
	binary.Write(buffer, binary.LittleEndian, header.size)
	binary.Write(buffer, binary.LittleEndian, header.hType)

	after := buffer.Len()
	return after - prev
}
