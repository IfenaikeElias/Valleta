package main

import (
	"net"
	"strings"
	"time"
    "strconv"
)

// write function to send response (generalize for different type of responses )
// implement different route handlers

// Processing response
func buildResponseheader(statusLine, connectionType, contentType, contentEncoding string, bodyLen int) []byte {
    var b strings.Builder
    b.Grow(256)

    b.WriteString(statusLine)
    b.WriteString("\r\n")

    if connectionType == "" {
        connectionType = "keep-alive"
    }
    b.WriteString("Connection: ")
    b.WriteString(connectionType)
    b.WriteString("\r\n")

    if contentType != "" {
        b.WriteString("Content-Type: ")
        b.WriteString(contentType)
        b.WriteString("\r\n")
    }

    if contentEncoding != "" {
        b.WriteString("Content-Encoding: ")
        b.WriteString(contentEncoding)
        b.WriteString("\r\n")
    }

    b.WriteString("Content-Length: ")
    b.WriteString(strconv.Itoa(bodyLen))
    b.WriteString("\r\n")

    b.WriteString("Date: ")
    b.WriteString(time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"))
    b.WriteString("\r\n")

    b.WriteString("Server: ")
    b.WriteString(name)
    b.WriteString("\r\n\r\n")

    return []byte(b.String())
}

func (s *Server) SendResponse(statusCode string, conn net.Conn, body string, headercontents ...string) {
    var (
        connectionType  string
        contentType     string
        contentEncoding string
    )

    switch len(headercontents) {
    case 3:
        contentEncoding = headercontents[2]
        fallthrough
    case 2:
        contentType = headercontents[1]
        fallthrough
    case 1:
        connectionType = headercontents[0]
    }

    if contentEncoding != "" && contentEncoding != "gzip" {
        contentEncoding = ""
    }

    bodyBytes := []byte(body)
    headerBytes := buildResponseheader("HTTP/1.1 "+statusCode, connectionType, contentType, contentEncoding, len(bodyBytes))

    conn.Write(append(headerBytes, bodyBytes...))

    if strings.ToLower(connectionType) == "close" {
        conn.Close()
    }
}
