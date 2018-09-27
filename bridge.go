package vishnu

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Bridge struct {
	port       int
	serverAddr string
}

func NewBridge(port int) *Bridge {
	return &Bridge{
		port: port,
	}
}

func (b *Bridge) Run() {
	server, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", b.port))
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go b.processConnection(conn)
	}
}

func (b *Bridge) processConnection(conn net.Conn) {
	isServer := false
	defer b.closeConnection(conn, &isServer)

	addr := conn.RemoteAddr().String()
	log.Printf("accept connection: %s", addr)
	fmt.Fprintln(conn, bridgeTag)

loop:
	for {
		d, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println(err)
			break
		}
		// 去除最后的换行
		d = d[:len(d)-1]

		switch d {
		case clientTag:
			log.Printf("client request to connect server: %s", addr)
			if b.serverAddr == "" {
				fmt.Fprintln(conn, noServerTag)
			} else {
				fmt.Fprintln(conn, b.serverAddr)
			}
		case serverTag:
			if b.serverAddr != "" {
				fmt.Fprintf(conn, "server already registerd at: %s", b.serverAddr)
				break loop
			}
			b.serverAddr = addr
			isServer = true
			log.Printf("server registered to bridge: %s", addr)
		default:
			log.Printf("%s", d)
		}
	}
}

func (b *Bridge) closeConnection(conn net.Conn, isServer *bool) {
	if *isServer {
		b.serverAddr = ""
	}
	conn.Close()
}
