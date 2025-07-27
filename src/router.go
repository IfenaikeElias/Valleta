package main

import (
	"fmt"
	"net"
	"os"
	"strings"
    "sync"
)
/*
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
*/

func (s *Server) route(conn net.Conn, req *HTTP_Request) {
    // fast path: root
    if req.url_path == "/" {
        if h, ok := s.Router["/"]; ok {
            h.Respond(conn, req, s)
            return
        }
    }

    // peel first segment without allocating a slice for every component
    trimmed := strings.TrimPrefix(req.url_path, "/")
    end := strings.IndexByte(trimmed, '/')
    if end == -1 {
        end = len(trimmed)
    }
    route_key := "/" + trimmed[:end]

    if handler, ok := s.Router[route_key]; ok {
        fmt.Println("Accepted connection from", conn.RemoteAddr())
        handler.Respond(conn, req, s)
        return
    }

    //  404 fallback 
    // Cache the not‑found page in memory so we don’t hit the filesystem
    // on every unknown route. A real server might embed the asset with
    // go:embed, but we’ll keep it simple with a package‑level var.
    not_found_once.Do(func() {
        data, err := os.ReadFile("web/Pagenotfound.html")
        if err != nil {
            not_found_page = []byte("<h1>404 page not found</h1>")
            return
        }
        not_found_page = data
    })

    body, encoding := compress_if_allowed(req, not_found_page)
    s.SendResponse("404 Not Found", conn, body, default_connection(req), "text/html", encoding)
}

var (
    not_found_page []byte
    not_found_once sync.Once
)

