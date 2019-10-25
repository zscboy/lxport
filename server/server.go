package server

import (
	"net"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"net/http"
)

var upgrader = websocket.Upgrader{} // use default options

func wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	var portStr = r.URL.Query().Get("port")
	if portStr == "" {
		log.Println("need port!")
		return
	}

	tcp, err := net.Dial("tcp", "127.0.0.1:"+portStr)
	if err != nil {
		log.Printf("dial to 127.0.0.1:%s failed: %v", portStr, err)
		return
	}

	defer tcp.Close()

	recvBuf := make([]byte, 4096)
	go func() {
		n, err := tcp.Read(recvBuf)
		if err != nil {
			log.Println("read from tcp failed:", err)
			c.Close()

			return
		}

		if n == 0 {
			log.Println("read from tcp got 0 bytes")
			c.Close()

			return
		}

		log.Println("tcp recv message, len:", n)
		tcp.Write(recvBuf[:n])
	}()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("websocket read failed:", err)
			tcp.Close()

			break
		}

		log.Println("websocket recv message, len:", len(message))
		err = writeAll(message, tcp)
		if err != nil {
			log.Println("write all to tcp failed:", err)
			break
		}
	}
}

func writeAll(buf []byte, nc net.Conn) error {
	wrote := 0
	l := len(buf)
	for {
		n, err := nc.Write(buf[wrote:])
		if err != nil {
			return err
		}

		wrote = wrote + n
		if wrote == l {
			break
		}
	}

	return nil
}

// CreateHTTPServer start http server
func CreateHTTPServer(listenAddr string, wsPath string) {
	http.HandleFunc(wsPath, wsHandler)
	log.Printf("server listen at:%s, path:%s", listenAddr, wsPath)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
