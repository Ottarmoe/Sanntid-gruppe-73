package main

import (
    "fmt"
    "net"
    "time"
)

func main() {
    // Bind to all interfaces, port 30000 to receive server broadcast
    addr := net.UDPAddr{
        IP:   net.IPv4zero, // 0.0.0.0 = all interfaces
        Port: 30000,
    }

    conn, err := net.ListenUDP("udp", &addr)
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    fmt.Println("Listening for server broadcast on UDP port 30000...")

    buf := make([]byte, 1024)
    for {
        // Receive broadcast
        n, remoteAddr, err := conn.ReadFromUDP(buf)
        if err != nil {
            fmt.Println("Receive error:", err)
            continue
        }

        // Print broadcast message and sender IP
        msg := string(buf[:n])
        fmt.Printf("Broadcast from %s: %s\n", remoteAddr.IP.String(), msg)

        // Optional: wait a second before listening again
        time.Sleep(time.Second)
    }
}
