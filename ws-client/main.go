package main

import (
	"bufio"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func main() {
	// WebSocket 서버 주소
	u := url.URL{Scheme: "ws", Host: "localhost:30000", Path: "/"}

	// HTTP 클라이언트를 사용하여 핸드셰이크 요청 생성
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v\n", err)
	}

	// WebSocket 핸드셰이크 헤더 설정
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Sec-WebSocket-Version", "13")
	req.Header.Set("Sec-WebSocket-Key", generateSecWebSocketKey())

	// WebSocket 서버에 연결

	conn, _, _, err := ws.Dial(context.TODO(), req.URL.String())
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

		// 메시지 전송
		err = wsutil.WriteClientText(conn, []byte(text))
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
