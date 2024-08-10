package socket

import (
	"sync/atomic"

	"github.com/JoTaeYang/study-go/pkg/lock/queue"
	"github.com/gobwas/ws"
	"github.com/panjf2000/gnet"
)

type SESSION_MODE int32

const (
	MODE_NONE SESSION_MODE = iota
	MODE_GAME
)

type Session struct {
	gnet.Conn
	recvBuffer    *RingBuffer
	completeRecvQ *queue.Queue[[]byte]
	sendQ         *queue.Queue[*ws.Frame]
	idx           int32
	mode          SESSION_MODE
	sendFlag      atomic.Int32
}
