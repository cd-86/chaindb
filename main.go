package main

import (
	"fmt"
	"net"
	"time"
)

func s() {
	addr := net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: 8888,
	}
	conn, err := net.DialUDP("udp", nil, &addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for {
		conn.Write([]byte("hello, this is a broadcast"))
		time.Sleep(time.Second * 5)
	}

}
func r() {
	addr := net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 8888,
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, senderAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		fmt.Printf("From %s: %s\n", senderAddr.IP.String(), string(buf[:n]))
	}
}
func main() {
	go s()
	go r()
	time.Sleep(10 * time.Second)
}
