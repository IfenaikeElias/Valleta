package main 

import (
	"io"
	"net"
	"os"
	"strings"
)



// routing : given a URL path, figure out which function (handler) shoudl process it

// this means any struct that implements Respond is a route and can be added to routes
type Handler interface {
	Respond(conn net.Conn, req *HTTP_Request, s *Server)
}

// created handler structs that will implement handler interface
type rootHandler struct{}
type echoHandler struct{}
type fileHandler struct{}

// these structs are routes, we map every url path to handlers
var routes = map[string]Handler{
	"/":      rootHandler{},
	"/echo":  echoHandler{},
	"/files": fileHandler{},
}


// methods for handling routes
func (h rootHandler) Respond(conn net.Conn, req *HTTP_Request, s *Server) {
	body := CompressResponsebody("Hey This is root")
	s.SendResponse("200 OK", conn, body, req.header_map["Connection"], "text/plain", req.header_map["Accept-Encoding"])
}

func (e echoHandler) Respond(conn net.Conn, req *HTTP_Request, s *Server) {
	msg := CompressResponsebody(strings.TrimPrefix(req.url_path, "/echo/"))
	s.SendResponse("200 OK", conn, msg, req.header_map["Connection"], "text/plain", req.header_map["Accept-Encoding"])
}

func (f fileHandler) Respond(conn net.Conn, req *HTTP_Request, s *Server) {
	var fileName, filePath string
	if strings.HasPrefix(req.url_path, "/files/") {
		fileName = strings.TrimPrefix(req.url_path, "/files/")
		filePath = s.Directory + "/" + fileName

		switch req.method {
		case "GET":
			fh, err := os.Open(filePath)
			if err != nil {
				s.SendResponse("404 Not Found", conn, "file you requested does not exist", req.header_map["Connection"], "text/plain", "")
				return
			}
			defer fh.Close()

			data, err := io.ReadAll(fh)
			if err != nil {
				s.SendResponse("500 Internal Server Error", conn, "cannot read file", req.header_map["Connection"], "text/plain", "")
				return
			}

			var contentType string
			if strings.HasSuffix(filePath, ".txt") {
				contentType = "text/plain"
			}

			if strings.HasSuffix(filePath, ".gz") {
				contentType = "application/octet-stream"
				// Send raw gzipped data with correct encoding
				s.SendResponse("200 OK", conn, CompressResponsebody(string(data)), req.header_map["Connection"], contentType, "gzip")
			} else {
				s.SendResponse("200 OK", conn, CompressResponsebody(string(data)), req.header_map["Connection"], contentType, "")
			}

		case "POST":
			// Ensure directory exists
			err := os.MkdirAll(s.Directory, 0755)
			if err != nil {
				s.SendResponse("500 Internal Server Error", conn, "failed to create directory", req.header_map["Connection"], "text/plain", "")
				return
			}
			fh, err := os.Create(filePath)
			if err != nil {
				s.SendResponse("500 Internal Server Error", conn, "cannot create file", req.header_map["Connection"], "text/plain", "")
				return
			}
			defer fh.Close()

			_, err = fh.Write([]byte(req.body))
			if err != nil {
				s.SendResponse("500 Internal Server Error", conn, "failed to write file", req.header_map["Connection"], "text/plain", "")
				return
			}
		}

		s.SendResponse("200 OK", conn, "created file successfully!", req.header_map["Connection"], "text/plain", "")
	}
}


