package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)


type HTTP_Request struct {
	method       string
	url_path     string
	http_version string
	header_map   map[string]string
	body         string
}

// Parse request
func ParseRequest(conn net.Conn, Request *HTTP_Request) {

	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil && err != io.EOF {
		fmt.Println("End of current request data:", err)
	}

	request := string(buff[:n])
	fmt.Println("RAW REQUEST:\n", request)
	lines := strings.Split(request, CRLF)

	if len(lines) == 0 {
		fmt.Println("400 Bad Request: Empty request")
		return
	}

	request_line := strings.Split(lines[0], " ")
	if len(request_line) < 3 {
		fmt.Println("400 Bad Request: Malformed request line")
		return
	}

	// Parse request target line
	Request.url_path = request_line[1]
	Request.http_version = request_line[2]
	Request.method = request_line[0]
	Request.header_map = make(map[string]string)

	var j, i int // j to keep track of where body starts and i for loop building header_map
	for i = 1; i < len(lines); i++ {
		line := lines[i]
		if line == "" {
			j++ // keep moving till we get to start of body
			break
		}

		key_val := strings.SplitN(line, ":", 2)
		if len(key_val) == 2 {
			if strings.Contains(strings.TrimSpace(key_val[1]), "gzip") {
				Request.header_map[strings.TrimSpace(key_val[0])] = "gzip"
			}
			Request.header_map[strings.TrimSpace(key_val[0])] = strings.TrimSpace(key_val[1])
		}
	}

	if j < len(lines) {
		Request.body = strings.Join(lines[i:], CRLF)
	} else {
		Request.body = ""
	}
	fmt.Println(Request.body)
	if err != nil {
		fmt.Println("400 Bad Request:", err)
	}
}
