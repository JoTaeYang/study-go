package socket

import (
	"bytes"
	"sync/atomic"

	"github.com/JoTaeYang/study-go/pkg/lock/queue"
	"github.com/JoTaeYang/study-go/pkg/ringbuffer"
	"github.com/panjf2000/gnet"
)

type SESSION_MODE int32

const (
	MODE_NONE SESSION_MODE = iota
	MODE_GAME
)

type Session struct {
	gnet.Conn
	recvBuffer    *ringbuffer.Buffer
	completeRecvQ *queue.Queue[[]byte]
	sendQ         *queue.Queue[*bytes.Buffer]
	idx           int32
	mode          SESSION_MODE
	sendFlag      atomic.Int32
}
