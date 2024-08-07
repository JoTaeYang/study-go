package yws

import (
	"bytes"
	"encoding/binary"
	"sync/atomic"

	"github.com/JoTaeYang/study-go/pkg/lockfree/queue"
	"github.com/gobwas/ws"
	"github.com/panjf2000/gnet"
)

type SESSION_MODE int32

const (
	MODE_NONE SESSION_MODE = iota
	MODE_GAME
)

type WebSocketConn struct {
	gnet.Conn
	completeRecvQ *queue.Queue[[]byte]
	sendQ         *queue.Queue[*ws.Frame]
	h             *ws.Frame
	Upgraded      bool
	idx           int32
	mode          SESSION_MODE
	bufIdx        atomic.Int32
	sendFlag      atomic.Int32
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

	//msg := wsc.makeWriteHeader(fr.Header)

	//ws.NewFrame으로 만들어서 패킷 보내주면 된다.
	//지금은 에코 서버 만들어서 이런 것.
	wsc.h.Header.Length = 0
	return fr.Payload
}
