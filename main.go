package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {

	store := NewStore()

	if err := startServer(store); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func startServer(store *Store) error {
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
		go handleConnection(conn, store)
	}
}

func handleConnection(conn net.Conn, store *Store) {
	defer conn.Close()
	log.Println("connection from", conn.RemoteAddr())

	// Read lines from the connection
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		log.Println("line: ", line, "conn:", conn.RemoteAddr())
		msg, err := processLine(line, store)

		if err != nil {
			fmt.Fprintln(conn, "ERR "+err.Error())
		} else if len(msg) > 0 {
			fmt.Fprintln(conn, msg)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("conn error:", err)
		return
	}
}

func validateParts(parts []string, expectedParts uint8) error {
	if len(parts) != int(expectedParts) {
		return fmt.Errorf("wrong number of arguments for '%v'", parts[0])
	}
	return nil
}

func processLine(line string, store *Store) (string, error) {
	parts := strings.Fields(line)
	log.Println("parts: ", parts)
	if len(parts) == 0 {
		return "", nil
	}

	cmd := strings.ToUpper(parts[0])

	switch cmd {
	case "SET":
		err := validateParts(parts, 3)
		if err != nil {
			return "", err
		}
		key, val := parts[1], parts[2]
		store.Set(key, val)
		return "OK", nil
	case "GET":
		err := validateParts(parts, 2)
		if err != nil {
			return "", err
		}
		key := parts[1]
		val, ok := store.Get(key)
		if !ok {
			return "(nil)", nil
		}
		return val, nil
	case "DEL":
		err := validateParts(parts, 2)
		if err != nil {
			return "", err
		}
		key := parts[1]
		if store.Delete(key) {
			return "1", nil
		} else {
			return "0", nil
		}
	case "SUBSCRIBE":
		err := validateParts(parts, 2)
		if err != nil {
			return "", err
		}
		return "Did receive SUBSCRIBE cmd", nil
	case "PUBLISH":
		err := validateParts(parts, 3)
		if err != nil {
			return "", err
		}
		return "Did receive PUBLISH cmd", nil
	default:
		return "", fmt.Errorf("unknown command: %v", cmd)
	}
}
