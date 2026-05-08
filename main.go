package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	if err := startServer(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func startServer() error {
	// Open a port
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		return err
	}
	defer listener.Close()
	log.Println("listening on :6379")

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Println("connection from", conn.RemoteAddr())

	// Read lines from the connection
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)
		// Echo the line back
		fmt.Fprintln(conn, line)
	}

	if err := scanner.Err(); err != nil {
		log.Println("conn error:", err)
		return
	}
}
