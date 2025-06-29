package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

const CRLF = "\r\n"
const name = "myGoHTTP/1.0"

type Server struct {
	Listener  net.Listener // Network listeer
	Router    map[string]Handler
	Directory string
}

// struct to hold Processed HTTP request


func (s *Server) ListenforConn() {
	// Bind to a port so server listens for incomming connections
	// net.listen returns a listener value that implements the net.Listener interface
	
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Failed to bind to port 8080")
		os.Exit(1)
	}
	s.Listener = ln
}

func (s *Server) AcceptConn() (net.Conn, error) {
	// listener.Accepts returns our network connection which implements the net.Conn interface and an error
	conn, err := s.Listener.Accept()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (s *Server) CloseConn() {
	// function to close connection
	err := s.Listener.Close()
	if err != nil {
		fmt.Println("Failed to Close connection", err.Error())
	}
}





var directory = flag.String("directory", "", "Path to file directory")

func main() {
	flag.Parse()
	s := Server{
		Router: routes,
		Directory: *directory,
	}
	defer s.CloseConn()
	s.ListenforConn()
	for {
		conn, err := s.AcceptConn()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}
		// net.Conn interface implements the RemoteAdrr() method which return the remote network address if known
		fmt.Println("Accepted Connection from: ", conn.RemoteAddr())
		go func() {
			for {
				var req HTTP_Request
				ParseRequest(conn, &req)
				s.route(conn, &req)
				connHeader := req.header_map["Connection"]
				if connHeader != "keep-alive" {
					conn.Close()
					return
				}
			}
		}()
	}
}
