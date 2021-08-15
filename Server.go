//Monika Sharma
//L20504237
//Programming Lab #2

package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func isError(err error) { //Handles if any error
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR :", err.Error())
		os.Exit(1)
	}
}

var n_forks int
var n_philosophers int

func request(philo_string string, send chan string) { // Handles Request coming for Fork
	no, err := strconv.Atoi(philo_string)
	_ = err
	f_1 := no
	f_2 := (no + 1) % 5
	if lock_on_forks[f_1] == 0 && lock_on_forks[f_2] == 0 {
		lock_on_forks[f_1] = 1
		lock_on_forks[f_2] = 1
		send <- philo_string + "-accepted," + ip_addrs[f_1] + "," + ip_addrs[f_2]
	} else {
		send <- philo_string + "-declined"
	}
}

var lock_on_forks [5]int
var ip_addrs [5]string
var completed int
var ip_addrs_arr [10]*net.UDPAddr

func controlRegister(conn *net.UDPConn) { //This is to control the ipAddress for UDP connection
	var buf [512]byte
	n, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}
	if string(buf[0:n]) == string("P") && (n_philosophers < 5) {
		ip_addrs_arr[n_philosophers] = addr

		conn.WriteToUDP([]byte(strconv.Itoa(n_philosophers)), addr)
		n_philosophers = n_philosophers + 1
	} else if strings.Split(string(buf[0:n]), ",")[0] == "F" && n_forks < 5 {
		ip_addrs_arr[5+n_forks] = addr
		ip_addrs[n_forks] = strings.Split(string(buf[0:n]), ",")[1]
		conn.WriteToUDP([]byte(strconv.Itoa(n_forks)), addr)
		n_forks += 1
	} else {
		conn.WriteToUDP([]byte("-1"), addr)
	}
}

func printStatus(no_string string, status_string string) { // This has all the Printing Status of Philosophers during the operation
	currentTime := time.Now().Format("Jan 06 15:04:05")
	print_string := currentTime + "\t\t"
	no, err := strconv.Atoi(no_string)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR :", err.Error())
		os.Exit(1)
	}
	status, err := strconv.Atoi(status_string)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR :", err.Error())
		os.Exit(1)
	}
	if status != 0 {
		for i := 0; i < 5; i++ {
			if no == i {
				switch status {
				case 0:
					print_string += "Thinking\t\t"
				case 1:
					print_string += "Waiting \t\t"
				case 2:
					print_string += "Eating  \t\t"
				case 3:
					print_string += "Finished\t\t"
				}
			} else {
				print_string += "........\t\t"
			}
		}
	} else {
		for i := 0; i < 5; i++ {
			print_string += "Thinking\t\t"
		}
	}
	fmt.Println(print_string)
}

func release(philo_string string, send chan string) { // When we have to release the fork
	no, err := strconv.Atoi(philo_string)
	_ = err

	f_1 := no
	f_2 := (no + 1) % 5

	lock_on_forks[f_1] = 0
	lock_on_forks[f_2] = 0
}

func controlSendChannel(send chan string, conn *net.UDPConn) { // This controls what we send through the channel
	for {
		buffer := <-send
		no, err := strconv.Atoi(strings.Split(buffer, "-")[0])
		message := strings.Split(buffer, "-")[1]
		_ = err
		conn.WriteToUDP([]byte(message), ip_addrs_arr[no])
	}
}

func controlRecvChannel(recv chan string, send chan string) { // This controls what we receive from the channel
	for {
		buffer := <-recv
		message := strings.Split(buffer, ",")

		c := message[0]
		if strings.Compare("print", c) == 0 {
			printStatus(message[1], message[2])
		} else if strings.Compare("request", c) == 0 {
			request(message[1], send)
		} else if strings.Compare("release", c) == 0 {
			release(message[1], send)
		} else if strings.Compare("completed", c) == 0 {
			completed = completed + 1
			printStatus(message[1], message[2])
			if completed == 5 {
				fmt.Println("------DINING PHILOSOPHER PROBLEM SOLVED------")
			}
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s port\n", os.Args[0])
		os.Exit(1)
	}
	service := ":" + os.Args[1]
	udpAddr, err := net.ResolveUDPAddr("udp", service)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR :", err.Error())
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR :", err.Error())
		os.Exit(1)
	}
	//The server is ready to accept requests.
	fmt.Println("The Server is Ready to Print in Port Number :" + service)
	n_forks = 0
	n_philosophers = 0
	for {
		controlRegister(conn)
		if n_philosophers == 5 && n_forks == 5 {
			break
		}
	}
	var i int
	for i = 0; i < 10; i++ {
		conn.WriteToUDP([]byte("S"), ip_addrs_arr[i])
	}
	for i = 0; i < 5; i++ {
		lock_on_forks[i] = 0
	}
	fmt.Println("Current time\t\tPhil #0\t\tPhil #1\t\tPhil #2\t\tPhil #3\t\tPhil #4\t\t")
	go printStatus("0", "0")
	var recv = make(chan string)
	var send = make(chan string)
	completed = 0
	go controlRecvChannel(recv, send)
	go controlSendChannel(send, conn)
	var buf [512]byte
	for {
		n, _, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			continue
		}
		recv <- string(buf[0:n])
	}
}
