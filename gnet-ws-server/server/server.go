package server

import (
	"github.com/JoTaeYang/study-go/packet/stgo"
	"github.com/JoTaeYang/study-go/pkg/yws"
)

type GnetWsServer struct {
	*yws.WsServer
}

func InitMsgProc(server *GnetWsServer) {
	server.SetMsgProc(func(session *yws.WebSocketConn, h *stgo.PacketHeader) error {
		switch h.Pid {
		case int32(stgo.PacketID_CS_CONNECT):
		case int32(stgo.PacketID_CS_CHAT_USER_MSG):
		case int32(stgo.PacketID_CS_PONG):
		}
		return nil
	})
}
