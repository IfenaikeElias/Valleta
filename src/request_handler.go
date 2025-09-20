package main

import (
	"fmt"
	"io"
	"net"
	"strings"
	"bufio"
	"strconv"
)

type HTTP_Request struct {
    method       string
    url_path     string
    http_version string
    header_map   map[string]string // all keys stored lowercase
    body         []byte
}

// ParseRequest streams a single HTTP/1.x request off conn.
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
        if line == "\r\n" { // blank line indicates body starts
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

