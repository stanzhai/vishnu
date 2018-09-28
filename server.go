package vishnu

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	bridge    string
	localAddr string
	port      int
}

func NewServer(bridge string, port int) *Server {
	return &Server{
		bridge: bridge,
		port:   port,
	}
}

func (s *Server) Run() {
	conn, err := net.Dial("tcp", s.bridge)
	if err != nil {
		log.Fatal(err)
	}

	go s.processConnection(conn)
	select {}
}

func (s *Server) processConnection(conn net.Conn) {
	defer conn.Close()

	d, _ := bufio.NewReader(conn).ReadString('\n')
	if !strings.Contains(d, bridgeTag) {
		log.Println("bridge is invalid!")
		return
	}
	_, err := fmt.Fprintln(conn, serverTag)
	if err == nil {
		log.Println("registered this server to bridge!")
	}

	for {
		d, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println(err)
			break
		}
		// 去除最后的换行
		d = d[:len(d)-1]
		if !strings.Contains(d, connectCmdTag) {
			continue
		}
		log.Println(d)
		clientAddr := strings.Split(d, " ")[1]
		go s.punch(clientAddr)
	}
}

func (s *Server) punch(clientAddr string) {

}
