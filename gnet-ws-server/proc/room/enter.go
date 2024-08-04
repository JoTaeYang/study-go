package room

import (
	"bytes"

	"github.com/JoTaeYang/study-go/packet"
	"github.com/JoTaeYang/study-go/packet/stgo"
	"github.com/JoTaeYang/study-go/pkg/yws"
	"github.com/gobwas/ws"
)

func Enter(session *yws.WebSocketConn, payload *[]byte) (fr *ws.Frame, err error) {
	payloadLen := len(*payload)
	buffer := bytes.Buffer{}
	buffer.Grow(16 + payloadLen)

	// create proto header
	header := stgo.PacketHeader{
		Pid:  int32(stgo.PacketID_SC_ECHO),
		Size: int32(payloadLen),
	}

	packet.HeaderToByte(&header, &buffer)

	_, err = buffer.Write(*payload)
	if err != nil {
		return
	}

	fr = &ws.Frame{}
	*fr = ws.NewFrame(ws.OpBinary, true, buffer.Bytes())
	return
}
