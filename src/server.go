package main


import (
	"fmt"
	"net"
	"os"
)

const CRLF = "\r\n"
const name = "myGoHTTP/1.0"


type Server struct {
    Listener  net.Listener
    Router    map[string]Handler
    Directory string
}

func (s *Server) ListenforConn(addr string) {
    ln, err := net.Listen("tcp", addr)
    if err != nil {
        fmt.Println("Failed to bind to", addr, "--", err)
        os.Exit(1)
    }
    s.Listener = ln
    fmt.Println("Listening on", addr)
}

func (s *Server) AcceptConn() (net.Conn, error) {
    return s.Listener.Accept()
}

func (s *Server) CloseConn() {
    if s.Listener != nil {
        _ = s.Listener.Close()
    }
}

