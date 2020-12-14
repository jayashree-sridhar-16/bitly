package main

import (
	"crypto/md5"
	"io"
	"fmt"
	"encoding/hex"
)

func encode_url(url string) string {
	h := md5.New()
	io.WriteString(h, url)	
	hash := h.Sum(nil)[:4]
	fmt.Printf("%x\n", hash)
	a := BytesToString(hash)
	return a
}

func BytesToString(data []byte) string {
	return hex.EncodeToString(data[:])
}