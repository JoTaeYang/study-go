package socket

import (
	"bytes"
	"log"
	"time"

	"github.com/JoTaeYang/study-go/pkg/lock/queue"
	"github.com/JoTaeYang/study-go/pkg/lock/stack"
	"github.com/JoTaeYang/study-go/pkg/ringbuffer"
	"github.com/gobwas/ws"
	"github.com/panjf2000/gnet"
)

type Handler interface {
	OnClientJoin(idx int32)
}

type Header struct {
	code  byte
	size  byte
	hType byte
}

type Server struct {
	*gnet.EventServer

	idx         *stack.Stack[int32]
	sessionList []*Session

	handler Handler
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
	s.idx = stack.NewStack[int32]()
	s.sessionList = make([]*Session, poolIdxLength)
	for i := int32(0); i < poolIdxLength; i++ {
		s.idx.Push(i)

		s.sessionList[i] = &Session{
			mode:          MODE_NONE,
			recvBuffer:    ringbuffer.NewBuffer(2048),
			completeRecvQ: queue.NewQueue[[]byte](),
			sendQ:         queue.NewQueue[*bytes.Buffer](),
		}
	}
}

func (s *Server) SetHandler(handler Handler) {
	s.handler = handler
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
					msgBuffer := make([]byte, 0, 1024)
					for i := int32(0); i < loopCnt; i++ {
						buf := v.sendQ.Dequeue()
						if buf == nil {
							v.Conn.Close()
							break
						}
						msgBuffer = append(msgBuffer, buf.Bytes()...)
					}
					v.AsyncWrite(msgBuffer)
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
	log.Printf("Socket server started on %s\n", srv.Addr.String())
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

	s.handler.OnClientJoin(idx)

	c.SetContext(wsc)
	return
}

/*
내가 Recv 중에 또 Recv가 들어올 수 있나?
*/

func (s *Server) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	wsc := c.Context().(*Session)

	transferred := len(frame)

	wsc.recvBuffer.Enqueue(&frame, int32(transferred))

	tmpFrame := make([]byte, transferred)

	copy(tmpFrame, frame)
	var outSize int32
	var err error
	for {
		if transferred <= 0 {
			break
		}
		msg := make([]byte, 0, 15)
		outSize, err = wsc.recvBuffer.Peek(&msg, 3)
		if err != nil {
			break
		}

		if msg[0] != 0x89 {
			break
		}

		if byte(wsc.recvBuffer.GetUseSize()) < msg[1]+3 {
			break
		}

		outSize, _ = wsc.recvBuffer.Dequeue(&msg, int32(msg[1]+3))

		wsc.completeRecvQ.Enqueue(msg)
		transferred -= int(outSize)
	}

	return
}

func (s *Server) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	wsc := c.Context().(*Session)
	log.Printf("Socket connection closed from %s\n", c.RemoteAddr().String())

	wsc.Conn = nil
	wsc.mode = MODE_NONE

	wsc.recvBuffer.Clear()
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

func (s *Server) SendPacket(idx int32, buffer *bytes.Buffer) {
	s.sessionList[idx].sendQ.Enqueue(buffer)
}
