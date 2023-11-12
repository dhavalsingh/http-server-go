package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	// buffer, err := io.ReadAll(conn)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)	
	}
	request := string(buffer)
	lines := strings.Split(request, "\r\n")
	start_line := lines[0]
	start_line_parts = strings.Fields(start_line)
	path := start_line_parts[1]

	var response string

	if path == "/"{
		response = "HTTP/1.1 200 OK\r\n\r\n"
	}
	else {
		response = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error sending msg: ", err.Error())
		os.Exit(1)
	}
}
