package main 

import (
	"io"
	"net"
	"os"
	"strings"
    "path/filepath"
    "mime"
)



// routing : given a URL path, figure out which function (handler) shoudl process it

// any struct that implements Respond is a route
// (kept identical to your interface name).
type Handler interface {
    Respond(conn net.Conn, req *HTTP_Request, s *Server)
}

type rootHandler struct{}
type echoHandler struct{}
type fileHandler struct{}

var routes = map[string]Handler{
    "/":      rootHandler{},
    "/echo":  echoHandler{},
    "/files": fileHandler{},
}


// compress_if_allowed returns a body string and the Content-Encoding value
// to be reported back to the client. It calls your existing
// CompressResponsebody helper only when the client advertises gzip support.
func compress_if_allowed(req *HTTP_Request, raw []byte) (string, string) {
    if strings.Contains(req.header_map["accept-encoding"], "gzip") {
        return CompressResponsebody(string(raw)), "gzip"
    }
    return string(raw), ""
}

func default_connection(req *HTTP_Request) string {
    if v, ok := req.header_map["connection"]; ok {
        return v
    }
    return "close"
}

//  handlers 

func (h rootHandler) Respond(conn net.Conn, req *HTTP_Request, s *Server) {
    body, encoding := compress_if_allowed(req, []byte("Hey This is root"))
    s.SendResponse("200 OK", conn, body, default_connection(req), "text/plain", encoding)
}

func (e echoHandler) Respond(conn net.Conn, req *HTTP_Request, s *Server) {
    msg := strings.TrimPrefix(req.url_path, "/echo/")
    body, encoding := compress_if_allowed(req, []byte(msg))
    s.SendResponse("200 OK", conn, body, default_connection(req), "text/plain", encoding)
}

func (f fileHandler) Respond(conn net.Conn, req *HTTP_Request, s *Server) {
    // validate prefix once and peel filename
    if !strings.HasPrefix(req.url_path, "/files/") {
        s.SendResponse("404 Not Found", conn, "invalid files route", default_connection(req), "text/plain", "")
        return
    }
    file_name := strings.TrimPrefix(req.url_path, "/files/")
    // stop path-traversal attempts like /files/../../etc/passwd
    file_name = filepath.Clean("/" + file_name)[1:]
    file_path := filepath.Join(s.Directory, file_name)

    switch req.method {
    case "GET":
        fh, err := os.Open(file_path)
        if err != nil {
            s.SendResponse("404 Not Found", conn, "file you requested does not exist", default_connection(req), "text/plain", "")
            return
        }
        defer fh.Close()

        data, err := io.ReadAll(fh)
        if err != nil {
            s.SendResponse("500 Internal Server Error", conn, "cannot read file", default_connection(req), "text/plain", "")
            return
        }

        ext := filepath.Ext(file_path)
        content_type := mime.TypeByExtension(ext)
        if content_type == "" {
            content_type = "application/octet-stream"
        }

        // If the file already ends with .gz, serve it with
        // Content-Encoding: gzip regardless of Accept-Encoding.
        if ext == ".gz" {
            s.SendResponse("200 OK", conn, string(data), default_connection(req), content_type, "gzip")
            return
        }

        body, encoding := compress_if_allowed(req, data)
        s.SendResponse("200 OK", conn, body, default_connection(req), content_type, encoding)

    case "POST":
        // ensure directory tree exists
        if err := os.MkdirAll(filepath.Dir(file_path), 0o755); err != nil {
            s.SendResponse("500 Internal Server Error", conn, "failed to create directory", default_connection(req), "text/plain", "")
            return
        }
        fh, err := os.Create(file_path)
        if err != nil {
            s.SendResponse("500 Internal Server Error", conn, "cannot create file", default_connection(req), "text/plain", "")
            return
        }
        defer fh.Close()

        if _, err := fh.Write(req.body); err != nil {
            s.SendResponse("500 Internal Server Error", conn, "failed to write file", default_connection(req), "text/plain", "")
            return
        }
        s.SendResponse("201 Created", conn, "created file successfully!", default_connection(req), "text/plain", "")

    default:
        s.SendResponse("405 Method Not Allowed", conn, "method not allowed", default_connection(req), "text/plain", "")
    }
}

