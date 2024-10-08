package yws

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/JoTaeYang/study-go/packet"
	"github.com/JoTaeYang/study-go/packet/stgo"
	"github.com/JoTaeYang/study-go/pkg/lockfree/lfstack"
	"github.com/JoTaeYang/study-go/pkg/lockfree/queue"
	"github.com/gobwas/ws"
	"github.com/panjf2000/gnet"
)

type WsServer struct {
	*gnet.EventServer

	idx         *lfstack.Stack[int32]
	sessionList []*WebSocketConn

	msgFroc func(idx int32, msg []byte) (*ws.Frame, error)
}

const (
	MaxHeaderSize = 14
	MinHeaderSize = 2
)

type OpCode byte

const (
	bit0 = 0x80
	bit1 = 0x40
	bit2 = 0x20
	bit3 = 0x10
	bit4 = 0x08
	bit5 = 0x04
	bit6 = 0x02
	bit7 = 0x01

	len7  = int64(125)
	len16 = int64(^(uint16(0)))
	len64 = int64(^(uint64(0)) >> 1)
)

/*
서버 세팅

@params poolIdxLength 세션 인덱스 보관 stack 의 사이즈
*/
func (s *WsServer) InitServer(poolIdxLength int32) {
	s.idx = lfstack.NewStack[int32]()
	s.sessionList = make([]*WebSocketConn, poolIdxLength)
	for i := int32(0); i < poolIdxLength; i++ {
		s.idx.Push(i)

		s.sessionList[i] = &WebSocketConn{
			mode:          MODE_NONE,
			completeRecvQ: queue.NewQueue[[]byte](),
			sendQ:         queue.NewQueue[*ws.Frame](),
		}
	}
}

func (s *WsServer) SetMsgProc(proc func(idx int32, msg []byte) (*ws.Frame, error)) {
	s.msgFroc = proc
}

func (s *WsServer) SendGo() {
	for {
		for _, v := range s.sessionList {
			if v.mode == MODE_GAME {
				loopCnt := v.sendQ.GetCount()
				if loopCnt > 0 {
					if loopCnt > 200 {
						loopCnt = 100
					}

					msgBuf := make([]byte, 0, 1024)
					for i := int32(0); i < loopCnt; i++ {
						msg := &ws.Frame{}
						v.sendQ.Dequeue(&msg)
						sendBuf := v.makeWriteHeader(msg.Header)

						sendBuf = append(sendBuf, msg.Payload...)

						msgBuf = append(msgBuf, sendBuf...)
					}
					v.AsyncWrite(msgBuf)
				}
			}
		}
		time.Sleep(time.Microsecond * 5)
	}
}

func (s *WsServer) GameGo() {
	// goroutine stop 처리 추가
	for {
		for _, v := range s.sessionList {
			if v.mode == MODE_GAME {
				loopCnt := v.completeRecvQ.GetCount()
				if loopCnt > 0 {
					if loopCnt > 200 {
						loopCnt = 100
					}
					for i := int32(0); i < loopCnt; i++ {
						msg := make([]byte, 0, 10)
						v.completeRecvQ.Dequeue(&msg)
						s.msgFroc(v.idx, msg)
					}
				}
			}
		}
		time.Sleep(time.Microsecond * 5)
	}
}

func (s *WsServer) upgrade(wsc *WebSocketConn, br *bufio.Reader, action *gnet.Action) (out *bytes.Buffer) {
	req, err := http.ReadRequest(br)
	if err != nil {
		log.Printf("Failed to read request: %v\n", err)
	}

	if req.ProtoMajor != 1 || req.ProtoMinor != 1 {
		log.Printf("Unsupported HTTP version: %s\n", req.Proto)
		*action = gnet.Close
		return
	}

	if req.Header.Get("Upgrade") != "websocket" || req.Header.Get("Connection") != "Upgrade" {
		log.Printf("Invalid Upgrade or Connection header: Upgrade=%s, Connection=%s\n",
			req.Header.Get("Upgrade"), req.Header.Get("Connection"))
		*action = gnet.Close
		return
	}

	var buf bytes.Buffer
	buf.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
	buf.WriteString("Upgrade: websocket\r\n")
	buf.WriteString("Connection: Upgrade\r\n")
	key := req.Header.Get("Sec-WebSocket-Key")
	acceptKey := generateAcceptKey(key)
	buf.WriteString("Sec-WebSocket-Accept: " + acceptKey + "\r\n")
	buf.WriteString("\r\n")
	out = &buf
	return
}

func (s *WsServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("WebSocket server started on %s\n", srv.Addr.String())
	return
}

func (s *WsServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Printf("New connection from %s\n", c.RemoteAddr().String())
	idx, check := s.idx.Pop()
	if !check {
		return
	}
	wsc := s.sessionList[idx]
	wsc.Conn = c
	wsc.idx = idx
	wsc.mode = MODE_GAME
	wsc.bufIdx.Store(0)

	c.SetContext(wsc)
	return
}

func (s *WsServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	wsc := c.Context().(*WebSocketConn)

	//Upgrade Check
	if !wsc.Upgraded {
		br := bufio.NewReader(bytes.NewReader(frame))

		// WebSocket 핸드셰이크 처리
		outBuf := s.upgrade(wsc, br, &action)

		out = outBuf.Bytes()

		wsc.Upgraded = true
		return
	}

	buf := bytes.NewBuffer(frame)

	msg := wsc.ReadBytes(buf)

	if msg != nil {
		header := &stgo.PacketHeader{}

		err := packet.ByteToHeader(&msg, header)
		if err != nil {
			action = gnet.Close
			return
		}

		wsc.completeRecvQ.Enqueue(msg)

		// fr, err := s.msgFroc(wsc, msg)
		// if err != nil {
		// 	action = gnet.Close
		// 	return
		// }

		// sendBuf := wsc.makeWriteHeader(fr.Header)

		// sendBuf = append(sendBuf, fr.Payload...)

		// if sendBuf != nil {
		// 	out = sendBuf
		// }
	}
	return
}

func (s *WsServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	wsc := c.Context().(*WebSocketConn)
	log.Printf("WebSocket connection closed from %s\n", c.RemoteAddr().String())

	wsc.Conn = nil
	wsc.Upgraded = false
	wsc.h = nil
	wsc.mode = MODE_NONE

	loopCnt := wsc.completeRecvQ.GetCount()
	clearMsg := []byte{}
	dummy := &ws.Frame{}
	for i := int32(0); i < loopCnt; i++ {
		wsc.completeRecvQ.Dequeue(&clearMsg)
	}
	loopCnt = wsc.sendQ.GetCount()
	for i := int32(0); i < loopCnt; i++ {
		wsc.sendQ.Dequeue(&dummy)
	}

	s.idx.Push(wsc.idx)
	return
}

func (s *WsServer) SendPacket(idx int32, fr *ws.Frame) {
	s.sessionList[idx].sendQ.Enqueue(fr)
}

func generateAcceptKey(secWebSocketKey string) string {
	const magicString = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	hash := sha1.New()
	hash.Write([]byte(secWebSocketKey + magicString))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
