package server

import (
	"bytes"
	"encoding/binary"

	"github.com/gobwas/ws"
	"github.com/panjf2000/gnet"
)

type WebSocketConn struct {
	gnet.Conn
	Upgraded bool
	header   *ws.Frame
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

	if wsc.header == nil {
		wsc.header = &ws.Frame{Payload: []byte{}}
	}

	if wsc.header.Header.Length == 0 {
		wsHeader, err := ws.ReadHeader(buf)
		if err != nil {
			// 처리 추가하기
		}
		wsc.header.Header = wsHeader
		wsc.header.Payload = wsc.header.Payload[:0]
		return nil
	}

	if buf.Len() > 0 {
		wsc.header.Payload = append(wsc.header.Payload, buf.Bytes()...)
	}

	if wsc.header.Header.Length != int64(len(wsc.header.Payload)) {
		return nil
	}

	fr := ws.UnmaskFrameInPlace(*wsc.header)

	msg := wsc.makeWriteHeader(fr.Header)

	msg = append(msg, fr.Payload...)

	wsc.header.Header.Length = 0

	return msg
}
