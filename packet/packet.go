package packet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/JoTaeYang/study-go/packet/stgo"
)

func ByteToHeader(buffer *[]byte, header *stgo.PacketHeader) error {
	if header == nil {
		return errors.New("header is nil")
	}

	msgBuf := bytes.NewBuffer(*buffer)

	binary.Read(msgBuf, binary.LittleEndian, &header.Code)
	if header.Code != 0x89 {
		return errors.New("not equal code")
	}

	if msgBuf.Len() < 8 {
		return errors.New("length is less copy")
	}
	binary.Read(msgBuf, binary.LittleEndian, &header.Pid)
	binary.Read(msgBuf, binary.LittleEndian, &header.Size)
	return nil
}

func HeaderToByte(header *stgo.PacketHeader, buffer *bytes.Buffer) int {
	prev := buffer.Len()

	binary.Write(buffer, binary.LittleEndian, uint64(0x89))
	binary.Write(buffer, binary.LittleEndian, header.Pid)
	binary.Write(buffer, binary.LittleEndian, header.Size)
	after := buffer.Len()
	return after - prev
}
