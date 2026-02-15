package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	fixedPort = "34933"
	delimPort = "33546"
	bufSize   = 1024
)

func main() {
	server := "10.0.10.111" // change to your server IP
	mode := "delim"          // "fixed" or "delim"
	port := delimPort
	if mode == "fixed" {
		port = fixedPort
	}

	conn, err := net.Dial("tcp", server+":"+port)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Connected to", conn.RemoteAddr())

	// Start receiving in a separate goroutine
	go recvLoop(conn, mode)

	// Send "hello" every 2 seconds
	for {
		send(conn, "hello", mode)
		fmt.Println(">> hello")
		time.Sleep(2 * time.Second)
	}
}

func recvLoop(conn net.Conn, mode string) {
	buf := make([]byte, bufSize)
	var leftover string // for delim mode

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Disconnected:", err)
			return
		}
		data := string(buf[:n])

		if mode == "fixed" {
			fmt.Println("<<", strings.TrimRight(data, "\x00"))
		} else { // delim mode
			leftover += data
			for {
				idx := strings.IndexByte(leftover, 0) // find '\0'
				if idx == -1 {
					break
				}
				msg := leftover[:idx]
				fmt.Println("<<", msg)
				leftover = leftover[idx+1:]
			}
		}
	}
}

func send(conn net.Conn, msg, mode string) {
	if mode == "fixed" {
		buf := make([]byte, bufSize)
		copy(buf, msg)
		conn.Write(buf)
	} else {
		conn.Write([]byte(msg + "\x00"))
	}
}
