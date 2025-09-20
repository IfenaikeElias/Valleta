package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
)

var directoryFlag = flag.String("directory", "", "Path to file directory")

func main() {
    flag.Parse()

    server := Server{
        Router:    routes,
        Directory: *directoryFlag,
    }
    defer server.CloseConn()

    server.ListenforConn(":8080")

    for {
        conn, err := server.AcceptConn()
        if err != nil {
            fmt.Println("Failed to accept connection:", err)
            continue
        }

        go func(c net.Conn) {
            for {
                req, err := ParseRequest(c)
                if err != nil {
                    // client closed connection or parse error
                    return
                }

                server.route(c, req)

                // close if either side does not want keep-alive
                connHeader := strings.ToLower(req.header_map["connection"])
                if connHeader != "keep-alive" {
                    return
                }
            }
        }(conn)
    }
}

