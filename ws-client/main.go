package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func main() {
	// 서버 주소
	addr := "ws://localhost:30000"
	ctx := context.Background()
	// WebSocket 연결 생성
	conn, _, _, err := ws.DefaultDialer.Dial(ctx, addr)
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
