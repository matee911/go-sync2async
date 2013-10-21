package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

func handleRequest(connection net.Conn) {
	// ignore input, just output
	line, err := bufio.NewReader(connection).ReadBytes('\n')
	if err != nil {
		log.Fatal(err)
	}

	command := strings.ToUpper(string(line[:]))
	log.Print(command)

	sleep := rand.Intn(50) * 100
	time.Sleep(time.Duration(sleep) * time.Millisecond)

	io.Copy(connection, bytes.NewBufferString("PONG"))
	connection.Close()
}

func main() {
	sock, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		connection, err := sock.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleRequest(connection)
	}
}
