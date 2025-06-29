package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func (s *Server) route(conn net.Conn, req *HTTP_Request) {
	// get handler for url_path
	trimmed := strings.TrimPrefix(req.url_path, "/")
	parts := strings.SplitN(trimmed, "/", 2)

	routeKey := "/" + parts[0]
	handler, ok := s.Router[routeKey]

	if ok {
		fmt.Println("Accepted Connection from: ", conn.RemoteAddr())
		handler.Respond(conn, req, s)
	} else {
		errHtml, err := os.ReadFile("web/Pagenotfound.html")
		if err != nil {
			s.SendResponse("500 Internal Server Error", conn, "Failed to load error page.", req.header_map["Connection"], "text/plain", req.header_map["Accept-Encoding"])
			fmt.Println(err)
			return
		}
		s.SendResponse("404 Not Found", conn, string(errHtml), req.header_map["Connection"], "text/html", req.header_map["Accept-Encoding"])
	}
}