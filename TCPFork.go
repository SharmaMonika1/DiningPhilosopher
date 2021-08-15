//Monika Sharma
//L20504237
//Programming Lab #2

package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func controlClient(send chan string, conn net.Conn) { // Controls the flow of CLient
	go controlSendChannel(send, conn)
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		_ = err
		fmt.Println(string(buf[0:n]))
		if string(buf[0:n]) == "start" {
			send <- "1"
		} else if string(buf[0:n]) == "stop" {
			send <- "0"
		}
	}
}

func controlSendChannel(send chan string, conn net.Conn) { // Controls the sends through the channel
	for {
		m := <-send
		if m == "0" {
			break
		}
		time.Sleep(time.Second)
		conn.Write([]byte("eat"))
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s server_host:server_port fork_host:fork_port", os.Args[0])
		os.Exit(1)
	}
	host := os.Args[1]
	service := os.Args[2]
	udpAddr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR : ", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR : ", err.Error())
		os.Exit(1)
	}
	_, err = conn.Write([]byte("F," + service))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR : ", err.Error())
		os.Exit(1)
	}
	var buf [512]byte
	n, err := conn.Read(buf[0:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR : ", err.Error())
		os.Exit(1)
	}
	no := string(buf[0:n])
	if no == "-1" {
		fmt.Fprintf(os.Stderr, "Server Rejected\n")
		os.Exit(1)
	}
	fmt.Println("Fork-no:", no)
	fmt.Println("Waiting for Server")
	n, err = conn.Read(buf[0:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR : ", err.Error())
		os.Exit(1)
	}
	fmt.Println("started")
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR : ", err.Error())
		os.Exit(1)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR : ", err.Error())
		os.Exit(1)
	}
	fmt.Println("Listening started....", no)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		send := make(chan string)
		go controlClient(send, conn)
	}
	os.Exit(0)
}
