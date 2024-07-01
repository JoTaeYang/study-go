package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"log"
	"net"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/panjf2000/gnet"
)

type wsServer struct {
	*gnet.EventServer
}

func (s *wsServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("WebSocket server started on %s\n", srv.Addr.String())
	return
}

func (s *wsServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Printf("New connection from %s\n", c.RemoteAddr().String())
	c.SetContext(new(WebSocketConn))
	return
}

func (s *wsServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	wsc := c.Context().(*WebSocketConn)

	if !wsc.Upgraded {
		// WebSocket 핸드셰이크 처리
		br := bufio.NewReader(bytes.NewReader(frame))
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

		resp := http.Response{
			Status:        "101 Switching Protocols",
			StatusCode:    http.StatusSwitchingProtocols,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			Header:        make(http.Header),
			ContentLength: -1,
		}
		key := req.Header.Get("Sec-WebSocket-Key")
		acceptKey := generateAcceptKey(key)
		resp.Header.Set("Upgrade", "websocket")
		resp.Header.Set("Connection", "Upgrade")
		resp.Header.Set("Sec-WebSocket-Accept", acceptKey)

		var buf bytes.Buffer
		resp.Write(&buf)
		out = buf.Bytes()

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

	err = wsutil.WriteServerMessage(wsc.Conn, ws.OpBinary, msg)
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
	net.Conn
	Upgraded bool
}

func (wsc *WebSocketConn) Read(b []byte) (n int, err error) {
	return wsc.Conn.Read(b)
}

func (wsc *WebSocketConn) Write(b []byte) (n int, err error) {
	return wsc.Conn.Write(b)
}

func generateAcceptKey(secWebSocketKey string) string {
	const magicString = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	hash := sha1.New()
	hash.Write([]byte(secWebSocketKey + magicString))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func main() {
	ws := &wsServer{}
	// WebSocket 서버 시작
	log.Fatal(gnet.Serve(ws, "tcp://:30000", gnet.WithMulticore(true)))
}
