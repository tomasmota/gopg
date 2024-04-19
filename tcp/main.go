package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

func runServer(readyChan chan<- struct{}, msgChan chan<- int64) {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Println("server listening on 8080")

	readyChan <- struct{}{}

	conn, err := ln.Accept()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	command, err := reader.ReadByte()
	if err != nil {
		log.Fatalf("error parsing command: %v", err)
	}
	fmt.Printf("Command: %c\n", command)

	lengthByte, err := reader.ReadByte()
	if err != nil {
		log.Fatalf("error reading payload length: %v", err)
	}
	length := uint8(lengthByte)
	fmt.Printf("length %d\n", length)

	key := make([]byte, length)
	n, err := io.ReadFull(reader, key)
	if err != nil {
		log.Fatalf("error reading key: %v", err)
	}
	if n != int(length) {
		log.Fatalf("expected to read %b bytes. got=%d", length, n)
	}
	fmt.Printf("key %s\n", key)

	var value int64
	err = binary.Read(reader, binary.BigEndian, &value)
	if err != nil {
		log.Fatalf("error reading value: %v", err)
	}
	// fmt.Printf("bytes client:%b", value.Bytes())

	newline, err := reader.ReadByte()
	if err != nil {
		log.Fatalf("error reading newline: %v", err)
	}
	if newline != '\n' {
		log.Fatalf("expected \n. got=%s", string(newline))
	}

	msgChan <- value
}

func runClient() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("error connecting to server: %v", err)
	}
	defer conn.Close()

	command := 'w'
	key := "some-key"
	if len(key) > 255 {
		log.Fatalf("key must be less than 255 bytes long. got=%d", len(key))
	}
	value := int64(434475932849)
	valueBuf := new(bytes.Buffer)
	err = binary.Write(valueBuf, binary.BigEndian, value)
	if err != nil {
		log.Fatalf("error parsing value into binary: %v", err)
	}
	fmt.Printf("bytes client:%b", valueBuf.Bytes())

	var payload bytes.Buffer
	payload.WriteByte(byte(command))
	payload.WriteByte(uint8(len(key)))
	payload.Write([]byte(key))
	payload.Write(valueBuf.Bytes())
	payload.WriteRune('\n')
	nBytes, err := conn.Write(payload.Bytes())
	if err != nil {
		log.Fatalf("error writing to connection: %v", err)
	}
	fmt.Printf("sent command (%d bytes)\n", nBytes)
}

func main() {
	readyChan := make(chan struct{})
	valueChan := make(chan int64)
	go runServer(readyChan, valueChan)
	<-readyChan
	runClient()
	fmt.Printf("Received value: %d\n", <-valueChan)
}
