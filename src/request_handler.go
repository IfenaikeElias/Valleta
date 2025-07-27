package main

import (
	"fmt"
	"io"
	"net"
	"strings"
	"bufio"
	"strconv"
)

/*
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

*/


type HTTP_Request struct {
    method       string
    url_path     string
    http_version string
    header_map   map[string]string // all keys stored lowercase
    body         []byte
}

// ParseRequest streams a single HTTP/1.x request off conn.
// No arbitrary limits; returns detailed errors to the caller.
func ParseRequest(conn net.Conn) (*HTTP_Request, error) {
    br := bufio.NewReader(conn)

    req := &HTTP_Request{header_map: make(map[string]string, 16)}

    // request line 
    line, err := br.ReadString('\n')
    if err != nil {
        return nil, err
    }
    line = strings.TrimRight(line, "\r\n")
    parts := strings.Fields(line)
    if len(parts) != 3 {
        return nil, fmt.Errorf("malformed request line %q", line)
    }
    req.method, req.url_path, req.http_version = parts[0], parts[1], parts[2]

    //  headers 
    for {
        line, err = br.ReadString('\n')
        if err != nil {
            return nil, err
        }
        if line == "\r\n" { // blank line = body starts
            break
        }
        kv := strings.SplitN(line, ":", 2)
        if len(kv) != 2 {
            return nil, fmt.Errorf("malformed header %q", line)
        }
        key := strings.TrimSpace(kv[0])
        val := strings.TrimSpace(kv[1])
        req.header_map[key] = val
    }

    // body 
    if cl_str, ok := req.header_map["content-length"]; ok {
        cl, err := strconv.Atoi(cl_str)
        if err != nil {
            return nil, fmt.Errorf("invalid Content-Length: %w", err)
        }
        req.body = make([]byte, cl)
        if _, err := io.ReadFull(br, req.body); err != nil {
            return nil, err
        }
    }

    return req, nil
}

