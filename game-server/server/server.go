package server

import (
	"bytes"
	"encoding/binary"
	randv2 "math/rand/v2"

	"github.com/JoTaeYang/study-go/packet/fighter"
	"github.com/JoTaeYang/study-go/pkg/socket"
)

type GameServer struct {
	*socket.Server

	player []*Player
}

func NewGameServer() *GameServer {
	player := make([]*Player, 1000, 1000)
	for i := 0; i < 1000; i++ {
		player[i] = &Player{}
	}

	gameServer := &GameServer{
		Server: &socket.Server{},
		player: player,
	}

	gameServer.SetHandler(gameServer)

	return gameServer
}

func (s *GameServer) OnClientJoin(idx int32) {

	msgBuf := make([]byte, 0, 64)
	buffer := bytes.NewBuffer(msgBuf)
	// create player packet make + player enter
	s.CreatePlayer(idx, buffer)

	// create player packet
	s.SendPacket(idx, buffer)
}

func (s *GameServer) CreatePlayer(idx int32, buffer *bytes.Buffer) {
	x := int16(randv2.Int64N(400) + 100)
	y := int16(randv2.Int64N(400) + 100)

	hp := byte(100)

	direction := byte(fighter.MoveDir_RR)

	s.player[idx].InitPlayer(idx, x, y, direction, 0, hp)

	binary.Write(buffer, binary.LittleEndian, byte(0x89))
	binary.Write(buffer, binary.LittleEndian, byte(10))
	binary.Write(buffer, binary.LittleEndian, byte(0))

	binary.Write(buffer, binary.LittleEndian, idx)
	binary.Write(buffer, binary.LittleEndian, direction)
	binary.Write(buffer, binary.LittleEndian, x)
	binary.Write(buffer, binary.LittleEndian, y)
	binary.Write(buffer, binary.LittleEndian, hp)
}

func InitMsgProc(server *GameServer) {
	server.SetMsgProc(nil)
}
