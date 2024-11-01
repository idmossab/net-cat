package main

import (
	"fmt"
	"log"
	"net"
)

type Message struct {
	from    string
	username string
	payload []byte
}

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgch      chan Message
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan Message, 10),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln
	go s.AcceptLoop()
	<-s.quitch
	close(s.msgch)
	return nil
}

func (s *Server) AcceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error :", err)
			continue
		}
		//fmt.Println("new connection to the server", conn.RemoteAddr())
		go s.ReadLoop(conn)
	}
}

func (s *Server) ReadLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	//for user
	n,err:=conn.Read(buf)
	if err!=nil{
		fmt.Println("read error:", err)
		return
	}
	usr:=string(buf[:n])
	fmt.Printf("User %s has connected",usr)
	//for msg
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}
		s.msgch <- Message{
			from:    conn.RemoteAddr().String(),
			username: usr,
			payload: buf[:n],
		}
		conn.Write([]byte("thank for your messag!\n"))
	}
}

func main() {
	server := NewServer(":3000")
	go func() {
		for msg := range server.msgch {
			fmt.Printf("received message from connection (%s):%s\n", msg.username, string(msg.payload))
		}
	}()
	log.Fatal(server.Start())
}
