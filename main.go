package main

import (
	"log"
	"net"
)

func main() {
	if err := startServer(); err != nil {
		log.Fatalf("Failed to start server, %v", err)
	}
}

func startServer() error {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		return err
	}
	defer listener.Close()
	log.Println("listening on :6379")
	
	conn, err := listener.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Println("connection from", conn.RemoteAddr())

	return nil
}