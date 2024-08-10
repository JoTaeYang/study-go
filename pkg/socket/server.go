package socket

import (
	"crypto/sha1"
	"encoding/base64"
	"log"
	"time"

	"github.com/JoTaeYang/study-go/pkg/lock/queue"
	"github.com/JoTaeYang/study-go/pkg/lockfree/lfstack"
	"github.com/gobwas/ws"
	"github.com/panjf2000/gnet"
)

type Header struct {
	code  byte
	size  byte
	hType byte
}

type Server struct {
	*gnet.EventServer

	idx         *lfstack.Stack[int32]
	sessionList []*Session

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
func (s *Server) InitServer(poolIdxLength int32) {
	s.idx = lfstack.NewStack[int32]()
	s.sessionList = make([]*Session, poolIdxLength)
	for i := int32(0); i < poolIdxLength; i++ {
		s.idx.Push(i)

		s.sessionList[i] = &Session{
			mode:          MODE_NONE,
			completeRecvQ: queue.NewQueue[[]byte](),
			sendQ:         queue.NewQueue[*ws.Frame](),
		}
	}
}

func (s *Server) SetMsgProc(proc func(idx int32, msg []byte) (*ws.Frame, error)) {
	s.msgFroc = proc
}

func (s *Server) SendGo() {
	for {
		for _, v := range s.sessionList {
			if v.mode == MODE_GAME {
				loopCnt := v.sendQ.GetCount()
				if loopCnt > 0 {
					if loopCnt > 200 {
						loopCnt = 100
					}

					v.AsyncWrite(msgBuf)
				}
			}
		}
		time.Sleep(time.Microsecond * 5)
	}
}

func (s *Server) GameGo() {
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
						//header dequeue
						msg := make([]byte, 0, 10)
						msg = v.completeRecvQ.Dequeue()
						s.msgFroc(v.idx, msg)
					}
				}
			}
		}
		time.Sleep(time.Microsecond * 5)
	}
}

func (s *Server) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("WebSocket server started on %s\n", srv.Addr.String())
	return
}

func (s *Server) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Printf("New connection from %s\n", c.RemoteAddr().String())
	idx, check := s.idx.Pop()
	if !check {
		return
	}
	wsc := s.sessionList[idx]
	wsc.Conn = c
	wsc.idx = idx
	wsc.mode = MODE_GAME

	c.SetContext(wsc)
	return
}

func (s *Server) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	wsc := c.Context().(*Session)

	transferred := len(frame)
	
	tmpFrame := make([]byte, transferred)

	copy(tmpFrame, frame)

	for {
		if transferred <= 0 {
			break
		}
		msg := &Header{}

		wsc.completeRecvQ.Enqueue(msg)
	}

	return
}

func (s *Server) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	wsc := c.Context().(*Session)
	log.Printf("WebSocket connection closed from %s\n", c.RemoteAddr().String())

	wsc.Conn = nil
	wsc.mode = MODE_NONE

	loopCnt := wsc.completeRecvQ.GetCount()

	for i := int32(0); i < loopCnt; i++ {
		wsc.completeRecvQ.Dequeue()
	}
	loopCnt = wsc.sendQ.GetCount()
	for i := int32(0); i < loopCnt; i++ {
		wsc.sendQ.Dequeue()
	}

	s.idx.Push(wsc.idx)
	return
}

func (s *Server) SendPacket(idx int32, fr *ws.Frame) {
	s.sessionList[idx].sendQ.Enqueue(fr)
}

func generateAcceptKey(secWebSocketKey string) string {
	const magicString = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	hash := sha1.New()
	hash.Write([]byte(secWebSocketKey + magicString))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
