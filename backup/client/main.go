package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// 서버에 연결
	conn, err := net.Dial("tcp", ":20000")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server. Type 'exit' to quit.")

	// 사용자 입력을 받아 서버에 전송
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter text: ")
		scanner.Scan()
		input := scanner.Text()

		if input == "exit" {
			break
		}

		// 서버에 메시지 전송
		fmt.Fprintf(conn, "%s\n", input)

		// 서버로부터 응답 수신 및 출력
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from server:", err)
			break
		}

		fmt.Println("Server response:", response)
	}

	fmt.Println("Connection closed.")
}
