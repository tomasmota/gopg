package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func runServer(msgChan chan<- string) {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Println("server listening on 8080")

	conn, err := ln.Accept()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var b strings.Builder
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		b.WriteString(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		if err == io.EOF {
			fmt.Println("connection closed by client")
		} else {
			log.Fatalf("error reading from connection: %v", err)
		}
	}
	msgChan <- b.String()
}

func runClient() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	nBytes, err := conn.Write([]byte("Hello\n from client\n"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("message sent, %d bytes\n", nBytes)
	// select {}
}

func main() {
	msgChan := make(chan string)
	go runServer(msgChan)
	go runClient()
	fmt.Printf("Received message: %s", <-msgChan)
}
