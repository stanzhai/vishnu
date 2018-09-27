package main

import (
	"io"
	"log"
	"net"

	reuse "github.com/libp2p/go-reuseport"
	"github.com/stanzhai/vishnu"
)

func main() {
	vishnu.Config.LoadConfig()
	config := vishnu.Config

	switch config.Type {
	case "client":
		vishnu.NewClient(config.Bridge, config.Port).Run()
	case "bridge":
		vishnu.NewBridge(config.Port).Run()
	case "server":
		vishnu.NewServer(config.Bridge, config.Port).Run()
	default:
		test()
	}
}

func test() {
	l, err := reuse.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			// Echo all incoming data.
			log.Printf("client addr: %s", c.RemoteAddr())
			//c2, e := net.Dial("tcp", "localhost:1235", "localhost:5000")
			c2, e := net.Dial("tcp", "localhost:5000")
			if e != nil {
				log.Fatal(e)
				return
			}
			log.Print("local addr: %s", c2.LocalAddr())
			_, e3 := reuse.Listen("tcp", c2.LocalAddr().String())
			if e3 != nil {
				log.Fatal(e3)
			}
			d := make(chan int, 1)
			go func() {
				io.Copy(c2, c)
				d <- 1
			}()
			go func() {
				io.Copy(c, c2)
				d <- 1
			}()
			// Shut down the connection.
			<-d
			<-d
			c2.Close()
			c.Close()
		}(conn)
	}
}
