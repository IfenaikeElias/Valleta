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
