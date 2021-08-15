//Monika Sharma
//L20504237
//Programming Lab #2

package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var conn *net.UDPConn

func isError(err error) { // Handles error that comes through
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR :", err.Error())
		os.Exit(1)
	}
}

func philosopher(no int, send chan string, recv chan string) { // All the philosophers performance is done here

	total_time := 9
	curr_status := 0
	eat_time := 0
	next_request_time := 0
	_ = next_request_time

	curr_status = 1
	send <- "print," + strconv.Itoa(no) + "," + strconv.Itoa(curr_status)
	var buf [512]byte
	for total_time > 0 {
		send <- "request," + strconv.Itoa(no)
		buffer := <-recv
		message := strings.Split(buffer, ",")
		if message[0] == "declined" {
			next_request_time = rand.Intn(total_time) + 1
			time.Sleep(time.Duration(next_request_time) * time.Second)
			continue
		}
		tcpAddr, err := net.ResolveTCPAddr("tcp4", message[1])
		if err != nil {
			send <- "release," + strconv.Itoa(no)
			next_request_time = rand.Intn(4) + 1
			send <- "print," + strconv.Itoa(no) + "," + strconv.Itoa(curr_status)
			time.Sleep(time.Duration(next_request_time) * time.Second)
			continue
		}
		conn1, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			send <- "release," + strconv.Itoa(no)
			next_request_time = rand.Intn(4) + 1
			send <- "print," + strconv.Itoa(no) + "," + strconv.Itoa(curr_status)
			time.Sleep(time.Duration(next_request_time) * time.Second)
			continue
		}
		tcpAddr, err = net.ResolveTCPAddr("tcp4", message[2])
		if err != nil {
			send <- "release," + strconv.Itoa(no)
			next_request_time = rand.Intn(4) + 1
			send <- "print," + strconv.Itoa(no) + "," + strconv.Itoa(curr_status)
			time.Sleep(time.Duration(next_request_time) * time.Second)
			continue
		}
		conn2, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			send <- "release," + strconv.Itoa(no)
			next_request_time = rand.Intn(4) + 1
			send <- "print," + strconv.Itoa(no) + "," + strconv.Itoa(curr_status)
			time.Sleep(time.Duration(next_request_time) * time.Second)
			continue
		}
		curr_status = 2
		send <- "print," + strconv.Itoa(no) + "," + strconv.Itoa(curr_status)
		eat_time = rand.Intn(total_time)
		if eat_time == 0 {
			eat_time = 1
		}
		i := 0
		for i < eat_time {
			conn1.Write([]byte("start"))
			n, err := conn1.Read(buf[0:])
			_ = err
			d := string(buf[0:n])
			conn2.Write([]byte("start"))
			n, err = conn2.Read(buf[0:])
			_ = err
			if d == "eat" && string(buf[0:n]) == "eat" {
				i = i + 1
			}
		}
		conn1.Write([]byte("stop"))
		conn2.Write([]byte("stop"))
		total_time = total_time - eat_time
		send <- "release," + strconv.Itoa(no)
		curr_status = 1
		send <- "print," + strconv.Itoa(no) + "," + strconv.Itoa(curr_status)
		if total_time >= 1 {
			next_request_time = rand.Intn(total_time) + 1
			time.Sleep(time.Duration(next_request_time) * time.Second)
		}
	}
	curr_status = 3
	send <- "completed," + strconv.Itoa(no) + "," + strconv.Itoa(curr_status)
	close(recv)
	for completed == 0 {
		next_request_time = rand.Intn(4) + 1
		time.Sleep(time.Duration(next_request_time) * time.Second)
	}
	close(send)
}

var completed int

func controlSendChannel(send chan string) { //Sent messages are handled in this function.
	for {
		buffer := <-send
		_, err := conn.Write([]byte(buffer))
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR :", err.Error())
			os.Exit(1)
		}
		if strings.Split(buffer, ",")[0] == "completed" {
			completed = 1
		}
	}
}

func controlRecvChannel(recv chan string) { //The received messages are handled in this function.
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR :", err.Error())
			os.Exit(1)
		}
		recv <- string(buf[0:n])
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port", os.Args[0])
		os.Exit(1)
	}
	service := os.Args[1]
	udpAddr, err := net.ResolveUDPAddr("udp", service)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR :", err.Error())
		os.Exit(1)
	}
	conn, err = net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR :", err.Error())
		os.Exit(1)
	}
	_, err = conn.Write([]byte("P"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR :", err.Error())
		os.Exit(1)
	}
	var buf [512]byte
	n, err := conn.Read(buf[0:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR :", err.Error())
		os.Exit(1)
	}
	no, _ := strconv.Atoi(string(buf[0:n]))
	if no == -1 {
		fmt.Fprintf(os.Stderr, "Server Rejected\n")
		os.Exit(1)
	}
	fmt.Println("Philo-no:", no)
	fmt.Println("Waiting for Server")
	n, err = conn.Read(buf[0:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR :", err.Error())
		os.Exit(1)
	}
	fmt.Println(string(buf[0:n]))
	fmt.Println("Starting...")
	var send = make(chan string)
	go controlSendChannel(send)
	var recv = make(chan string)
	go controlRecvChannel(recv)
	philosopher(no, send, recv)
	os.Exit(0)
}
