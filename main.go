package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
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
		log.Println("line: ", line, "conn:", conn.RemoteAddr())
		msg, err := processLine(line)

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

func processLine(line string) (string, error) {
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
		return "Did receive SET cmd", nil
	case "GET":
		err := validateParts(parts, 2)
		if err != nil {
			return "", err
		}
		return "Did receive GET cmd", nil
	case "DEL":
		err := validateParts(parts, 2)
		if err != nil {
			return "", err
		}
		return "Did receive DEL cmd", nil
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
