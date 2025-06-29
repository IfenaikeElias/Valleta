package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// write function to send response (generalize for different type of responses )
// implement different route handlers

// Processing response
func buildResponseheader(contentType, contentEncoding, connectionType, body string) map[string]string {
	header_map := make(map[string]string)

	if contentEncoding != "" {
		header_map["Content-encoding"] = contentEncoding
	}
	if contentType != "" {
		header_map["Content-type"] = contentType
	}
	if connectionType != "" {
		header_map["Connection"] = connectionType
	}
	header_map["Content-Length"] = fmt.Sprintf("%d", len([]byte(body)))
	header_map["Date"] = time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	header_map["Server"] = name

	return header_map
}

// SendResponse: After building response header, it construccts reponse and writes it our connection
// optional arguements -
// headercontents: headercontents[0] is connectionType , headercontents[1] is contentType, headercontents[2] is contentEncoding

// connectionType, contentType, contentEncoding,
func (s *Server) SendResponse(statusCode string, conn net.Conn, body string, headercontents ...string) {
	connectionType := "keep-alive" // Default connection type
	contentType := ""              // Default content type
	contentEncoding := ""          // Default contentEncoding
	if len(headercontents) > 0 {
		connectionType = headercontents[0]
	}
	if len(headercontents) > 1 {
		contentType = headercontents[1]
	}
	if len(headercontents) > 2 {
		contentEncoding = headercontents[2]
	}

	response_header_map := buildResponseheader(contentType, contentEncoding, connectionType, body)
	encoding, ok := response_header_map["Content-encoding"]
	if ok {
		if encoding != "gzip" {
			conn.Write([]byte("404 bad request: Unsupported encryption type"))
		}
	} else {
		var response strings.Builder
		response.WriteString(fmt.Sprintf("HTTP/1.1 %s\r\n", statusCode))

		for k, v := range response_header_map {
			response.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
		}
		response.WriteString("\r\n")
		response.WriteString(body)
		conn.Write([]byte(response.String()))
	}

}
