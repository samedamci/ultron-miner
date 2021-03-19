/*
Simple script to test SOCKS5 proxy connection.

Usage:
    go run test_proxy.go <addr:port>
*/

package main

import (
	"fmt"
	"golang.org/x/net/proxy"
	"net"
	"os"
	"time"
)

const POOL_ADDR string = "51.15.127.80:2811"

func check(e error) {
	if e != nil {
		if _, ok := e.(*net.OpError); ok {
			fmt.Println("proxy timeout")
		} else {
			fmt.Println(e)
			fmt.Printf("%T\n", e)
			os.Exit(1)
		}
	}
}

func main() {
	proxy_addr := os.Args[1]
	dialer, err := proxy.SOCKS5(
		"tcp", proxy_addr, nil, &net.Dialer{Timeout: 5 * time.Second},
	)
	check(err)
	conn, err := dialer.Dial("tcp", POOL_ADDR)
	check(err)
	fmt.Printf("connected to %s proxy\n", proxy_addr)
	buffer := make([]byte, 3)
	_, err = conn.Read(buffer)
	if err != nil {
		return
	}
	if len(buffer) >= 3 {
		fmt.Println("connected to pool, proxy is good")
	}
}
