package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/panjf2000/gnet/v2"
)

type wsServer struct {
	*gnet.BuiltinEventEngine

	addr      string
	multicore bool
	eng       gnet.Engine
	connected int64
}

func (s *wsServer) OnInitComplete(eng gnet.Engine) (action gnet.Action) {
	s.eng = eng
	return
}

func (s *wsServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Printf("New connection from %s\n", c.RemoteAddr().String())
	wsc := &WebSocketConn{Conn: c}
	c.SetContext(wsc)
	return
}

func (s *wsServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	wsc := c.Context().(*WebSocketConn)

	if !wsc.Upgraded {
		// WebSocket 핸드셰이크 처리
		buff, _ := c.Next(-1)
		br := bufio.NewReader(bytes.NewReader(buff))
		req, err := http.ReadRequest(br)
		if err != nil {
			log.Printf("Failed to read request: %v\n", err)
			action = gnet.Close
			return
		}

		if req.ProtoMajor != 1 || req.ProtoMinor != 1 {
			log.Printf("Unsupported HTTP version: %s\n", req.Proto)
			action = gnet.Close
			return
		}

		if req.Header.Get("Upgrade") != "websocket" || req.Header.Get("Connection") != "Upgrade" {
			log.Printf("Invalid Upgrade or Connection header: Upgrade=%s, Connection=%s\n",
				req.Header.Get("Upgrade"), req.Header.Get("Connection"))
			action = gnet.Close
			return
		}

		var buf bytes.Buffer
		buf.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
		buf.WriteString("Upgrade: websocket\r\n")
		buf.WriteString("Connection: Upgrade\r\n")
		key := req.Header.Get("Sec-WebSocket-Key")
		acceptKey := generateAcceptKey(key)
		buf.WriteString("Sec-WebSocket-Accept: " + acceptKey + "\r\n")
		buf.WriteString("\r\n")

		c.Write(buff)

		wsc.Upgraded = true

		return
	}

	// WebSocket 메시지를 에코
	msg, op, err := wsutil.ReadClientData(wsc)
	if err != nil {
		log.Printf("Failed to read WebSocket frame: %v\n", err)
		action = gnet.Close
		return
	}

	if op == ws.OpClose {

	}

	err = wsutil.WriteServerMessage(wsc, ws.OpBinary, msg)
	if err != nil {
		log.Printf("Failed to write WebSocket frame: %v\n", err)
		action = gnet.Close
	}

	return
}

func (s *wsServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	log.Printf("WebSocket connection closed from %s\n", c.RemoteAddr().String())
	return
}

type WebSocketConn struct {
	gnet.Conn
	Upgraded bool
}

func (wsc *WebSocketConn) Read(b []byte) (n int, err error) {

	data, err := wsc.Conn.Next(len(b))
	n = copy(b, data)
	if n == 0 {
		err = errors.New("no data read")
	}
	return
}

func (wsc *WebSocketConn) Write(b []byte) (n int, err error) {
	n, err = wsc.Conn.Write(b)
	if err != nil {
		return 0, err
	}
	return
}

func generateAcceptKey(secWebSocketKey string) string {
	const magicString = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	hash := sha1.New()
	hash.Write([]byte(secWebSocketKey + magicString))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func main() {

	// WebSocket 서버 시작
	var port int
	var multicore bool

	// Example command: go run main.go --port 8080 --multicore=true
	flag.IntVar(&port, "port", 30000, "server port")
	flag.BoolVar(&multicore, "multicore", true, "multicore")
	flag.Parse()

	ws := &wsServer{addr: fmt.Sprintf("tcp://127.0.0.1:%d", port), multicore: multicore}
	log.Println("server exits:", gnet.Run(ws, ws.addr, gnet.WithMulticore(multicore), gnet.WithReusePort(true), gnet.WithTicker(true)))
}
