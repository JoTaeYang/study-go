package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/JoTaeYang/study-go/packet"
	"github.com/JoTaeYang/study-go/packet/stgo"
	"github.com/JoTaeYang/study-go/pkg/convert"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func main() {
	// WebSocket 서버 주소
	u := url.URL{Scheme: "ws", Host: "localhost:30000", Path: "/"}

	// WebSocket 서버에 연결
	conn, _, _, err := ws.DefaultDialer.Dial(context.TODO(), u.String())
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket server: %v\n", err)
	}
	defer conn.Close()

	fmt.Println("Connected to WebSocket server")

	// 사용자 입력을 읽어 WebSocket 메시지로 전송
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message: ")
		text, _ := reader.ReadString('\n')

		header := stgo.PacketHeader{
			Pid:  int32(stgo.PacketID_CS_ECHO),
			Size: int32(len(text)),
		}

		buffer := bytes.Buffer{}

		packet.HeaderToByte(&header, &buffer)

		buffer.Write(convert.StringToBytes(text))

		// 메시지 전송
		err = wsutil.WriteClientText(conn, buffer.Bytes())
		if err != nil {
			log.Fatalf("Failed to send message: %v\n", err)
		}

		// 서버로부터 메시지 수신
		msg, op, err := wsutil.ReadServerData(conn)
		if err != nil {
			log.Fatalf("Failed to read message: %v\n", err)
		}

		if op == ws.OpClose {
			fmt.Println("Server closed the connection")
			break
		}

		fmt.Printf("Received: %s\n", string(msg))
	}
}

// Sec-WebSocket-Key 생성
func generateSecWebSocketKey() string {
	const magicString = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	key := "dGhlIHNhbXBsZSBub25jZQ==" // 임의의 Base64 인코딩 문자열
	hash := sha1.New()
	hash.Write([]byte(key + magicString))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
