package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

func CompressResponsebody(responseBody string) string {
	var buff bytes.Buffer
	gw := gzip.NewWriter(&buff)
	gw.Write([]byte(responseBody))
	gw.Close()
	compressed := buff.Bytes()
	return fmt.Sprintf("%x", compressed)
}

/* Optimized Version

var gzip_writer_pool = sync.Pool{
    New: func() interface{} {
        w, _ := gzip.NewWriterLevel(io.Discard, gzip.BestSpeed)
        return w
    },
}

// CompressResponsebody gzips the input string and returns the raw compressed
// frame as a string. Callers should set the "Content-Encoding: gzip" header.
func CompressResponsebody(responseBody string) string {
    var buff bytes.Buffer

    gw := gzip_writer_pool.Get().(*gzip.Writer)
    gw.Reset(&buff)

    _, _ = gw.Write([]byte(responseBody)) // bytes.Buffer never returns error
    gw.Close()

    gzip_writer_pool.Put(gw)

    return buff.String()
}


*/