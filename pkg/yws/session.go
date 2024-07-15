package yws

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/gobwas/ws"
	"github.com/panjf2000/gnet"
)

type WebSocketConn struct {
	gnet.Conn
	Upgraded bool
	h        *ws.Frame
	idx      int32
}

func (wsc *WebSocketConn) makeWriteHeader(h ws.Header) []byte {
	bts := make([]byte, MaxHeaderSize)

	if h.Fin {
		bts[0] |= bit0
	}
	bts[0] |= h.Rsv << 4
	bts[0] |= byte(h.OpCode)

	var n int
	switch {
	case h.Length <= len7:
		bts[1] = byte(h.Length)
		n = 2

	case h.Length <= len16:
		bts[1] = 126
		binary.BigEndian.PutUint16(bts[2:4], uint16(h.Length))
		n = 4

	case h.Length <= len64:
		bts[1] = 127
		binary.BigEndian.PutUint64(bts[2:10], uint64(h.Length))
		n = 10

	default:
		return nil
	}

	if h.Masked {
		bts[1] |= bit0
		n += copy(bts[n:], h.Mask[:])
	}

	return bts[:n]
}

func (wsc *WebSocketConn) ReadBytes(buf *bytes.Buffer) []byte {

	if wsc.h == nil {
		wsc.h = &ws.Frame{Payload: []byte{}}
	}

	if wsc.h.Header.Length == 0 {
		wsHeader, err := ws.ReadHeader(buf)
		if err != nil {
			// 처리 추가하기
		}
		wsc.h.Header = wsHeader
		wsc.h.Payload = wsc.h.Payload[:0]
		return nil
	}

	if buf.Len() > 0 {
		wsc.h.Payload = append(wsc.h.Payload, buf.Bytes()...)
	}

	if wsc.h.Header.Length != int64(len(wsc.h.Payload)) {
		return nil
	}

	fr := ws.UnmaskFrameInPlace(*wsc.h)

	msg := wsc.makeWriteHeader(fr.Header)

	//ws.NewFrame으로 만들어서 패킷 보내주면 된다.
	//지금은 에코 서버 만들어서 이런 것.

	log.Println(string(fr.Payload))

	msg = append(msg, fr.Payload...)

	wsc.h.Header.Length = 0

	return msg
}
