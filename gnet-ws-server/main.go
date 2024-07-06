package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"log"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/panjf2000/gnet"
)

type wsServer struct {
	*gnet.EventServer
}

const (
	MaxHeaderSize = 14
	MinHeaderSize = 2
)

type OpCode byte

const (
	bit0 = 0x80
	bit1 = 0x40
	bit2 = 0x20
	bit3 = 0x10
	bit4 = 0x08
	bit5 = 0x04
	bit6 = 0x02
	bit7 = 0x01

	len7  = int64(125)
	len16 = int64(^(uint16(0)))
	len64 = int64(^(uint64(0)) >> 1)
)

func (s *wsServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("WebSocket server started on %s\n", srv.Addr.String())
	return
}

func (s *wsServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Printf("New connection from %s\n", c.RemoteAddr().String())
	wsc := &WebSocketConn{Conn: c}
	c.SetContext(wsc)
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

		out = buf.Bytes()

		wsc.Upgraded = true
		return
	}

	buf := bytes.NewBuffer(frame)
	msg := wsc.ReadBytes(buf)

	//보낼 때도 write frame header 를 추가해줘야 한다.
	if msg != nil {
		out = msg
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
	readLen  int64
	header   *ws.Frame
}

func (wsc *WebSocketConn) Read(b []byte) (n int, err error) {
	wsc.Conn.ReadN(len(b))
	return
}

func (wsc *WebSocketConn) Write(b []byte) (n int, err error) {
	//wsc.Conn.Write(b)
	return
}

func (wsc *WebSocketConn) makeWriteHeader(h ws.Header) []byte {
	bts := make([]byte, MaxHeaderSize)

	if h.Fin {
		bts[0] |= bit0
	}
	bts[0] |= h.Rsv << 4
	bts[0] |= byte(h.OpCode)

	var n int
	switch {
	case h.Length <= len7:
		bts[1] = byte(h.Length)
		n = 2

	case h.Length <= len16:
		bts[1] = 126
		binary.BigEndian.PutUint16(bts[2:4], uint16(h.Length))
		n = 4

	case h.Length <= len64:
		bts[1] = 127
		binary.BigEndian.PutUint64(bts[2:10], uint64(h.Length))
		n = 10

	default:
		return nil
	}

	if h.Masked {
		bts[1] |= bit0
		n += copy(bts[n:], h.Mask[:])
	}

	return bts[:n]
}

func (wsc *WebSocketConn) ReadBytes(buf *bytes.Buffer) []byte {

	if wsc.header == nil {
		wsc.header = &ws.Frame{Payload: []byte{}}
	}

	if wsc.header.Header.Length == 0 {
		wsHeader, err := ws.ReadHeader(buf)
		if err != nil {

		}
		wsc.header.Header = wsHeader
		wsc.header.Payload = wsc.header.Payload[:0]
		return nil
	}

	if buf.Len() > 0 {
		wsc.header.Payload = append(wsc.header.Payload, buf.Bytes()...)
	}

	if wsc.header.Header.Length != int64(len(wsc.header.Payload)) {
		return nil
	}

	fr := ws.UnmaskFrameInPlace(*wsc.header)

	msg := wsc.makeWriteHeader(fr.Header)

	msg = append(msg, fr.Payload...)

	wsc.header.Header.Length = 0

	return msg
}

func main() {
	ws := &wsServer{}
	log.Fatal(gnet.Serve(ws, "tcp://:30000", gnet.WithMulticore(true)))
}

// type wsServer struct {
// 	*gnet.BuiltinEventEngine

// 	addr      string
// 	multicore bool
// 	eng       gnet.Engine
// 	connected int64
// }

// func (s *wsServer) OnInitComplete(eng gnet.Engine) (action gnet.Action) {
// 	s.eng = eng
// 	return
// }

// func (s *wsServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
// 	log.Printf("New connection from %s\n", c.RemoteAddr().String())
// 	wsc := &WebSocketConn{Conn: c}
// 	c.SetContext(wsc)
// 	return
// }

// func (s *wsServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
// 	wsc := c.Context().(*WebSocketConn)

// 	if !wsc.Upgraded {
// 		// WebSocket 핸드셰이크 처리
// 		buff, _ := c.Next(-1)
// 		br := bufio.NewReader(bytes.NewReader(buff))
// 		req, err := http.ReadRequest(br)
// 		if err != nil {
// 			log.Printf("Failed to read request: %v\n", err)
// 			action = gnet.Close
// 			return
// 		}

// 		if req.ProtoMajor != 1 || req.ProtoMinor != 1 {
// 			log.Printf("Unsupported HTTP version: %s\n", req.Proto)
// 			action = gnet.Close
// 			return
// 		}

// 		if req.Header.Get("Upgrade") != "websocket" || req.Header.Get("Connection") != "Upgrade" {
// 			log.Printf("Invalid Upgrade or Connection header: Upgrade=%s, Connection=%s\n",
// 				req.Header.Get("Upgrade"), req.Header.Get("Connection"))
// 			action = gnet.Close
// 			return
// 		}

// 		var buf bytes.Buffer
// 		buf.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
// 		buf.WriteString("Upgrade: websocket\r\n")
// 		buf.WriteString("Connection: Upgrade\r\n")
// 		key := req.Header.Get("Sec-WebSocket-Key")
// 		acceptKey := generateAcceptKey(key)
// 		buf.WriteString("Sec-WebSocket-Accept: " + acceptKey + "\r\n")
// 		buf.WriteString("\r\n")

// 		c.Write(buff)

// 		wsc.Upgraded = true

// 		return
// 	}

// 	// WebSocket 메시지를 에코
// 	msg, op, err := wsutil.ReadClientData(wsc)
// 	if err != nil {
// 		log.Printf("Failed to read WebSocket frame: %v\n", err)
// 		action = gnet.Close
// 		return
// 	}

// 	if op == ws.OpClose {

// 	}

// 	err = wsutil.WriteServerMessage(wsc, ws.OpBinary, msg)
// 	if err != nil {
// 		log.Printf("Failed to write WebSocket frame: %v\n", err)
// 		action = gnet.Close
// 	}

// 	return
// }

// func (s *wsServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
// 	log.Printf("WebSocket connection closed from %s\n", c.RemoteAddr().String())
// 	return
// }

// type WebSocketConn struct {
// 	gnet.Conn
// 	Upgraded bool
// }

// func (wsc *WebSocketConn) Read(b []byte) (n int, err error) {

// 	data, err := wsc.Conn.Next(len(b))
// 	n = copy(b, data)
// 	if n == 0 {
// 		err = errors.New("no data read")
// 	}
// 	return
// }

// func (wsc *WebSocketConn) Write(b []byte) (n int, err error) {
// 	n, err = wsc.Conn.Write(b)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return
// }

// func main() {

// 	// WebSocket 서버 시작
// 	var port int
// 	var multicore bool

// 	// Example command: go run main.go --port 8080 --multicore=true
// 	flag.IntVar(&port, "port", 30000, "server port")
// 	flag.BoolVar(&multicore, "multicore", true, "multicore")
// 	flag.Parse()

// 	ws := &wsServer{addr: fmt.Sprintf("tcp://127.0.0.1:%d", port), multicore: multicore}
// 	log.Println("server exits:", gnet.Run(ws, ws.addr, gnet.WithMulticore(multicore), gnet.WithReusePort(true), gnet.WithTicker(true)))
// }

func generateAcceptKey(secWebSocketKey string) string {
	const magicString = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	hash := sha1.New()
	hash.Write([]byte(secWebSocketKey + magicString))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
