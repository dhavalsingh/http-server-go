package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"flag"
	"io"
	"path/filepath"
)

func readFileContents (filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		// Handle the error, possibly a 404 if the file doesn't exist
		fmt.Println("Error accepting connection: ", err.Error())
		return "HTTP/1.1 404 Not Found\r\n\r\n"
	}
	defer file.Close()
	
	contents, err := io.ReadAll(file)
	if err != nil {
		// Handle the error
		fmt.Println("Error reading file: ", err.Error())
		return "HTTP/1.1 500 Internal Server Error\r\n\r\n"
	}
	return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(contents), contents)
}

func handleConnection (conn net.Conn, directory string){
	defer conn.Close()
	buffer := make([]byte, 1024)
	// buffer, err := io.ReadAll(conn)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		// os.Exit(1)	
	}
	request := string(buffer)
	lines := strings.Split(request, "\r\n")
	start_line := lines[0]
	// headers := lines[1]
	start_line_parts := strings.Fields(start_line)
	rMethod, rPath, rProtocol := start_line_parts[0], start_line_parts[1], start_line_parts[2]
	fmt.Printf("method=%s, path=%s, protocol=%s\n", rMethod, rPath, rProtocol)

	subRoute := strings.Split(rPath, "/")
	// ua := strings.Split(lines[2], " ")[1]

	var response string

	if rMethod == "POST" {
		fileName := strings.TrimPrefix(rPath, "/files/")
		filePath := filepath.Join(directory, fileName)

		// Find the start of the body
		bodyStartIndex := len(start_line) + 4 // Adjust this based on actual headers
		for i, line := range lines {
			if line == "" { // Empty line indicates end of headers
				bodyStartIndex += len(strings.Join(lines[:i], "\r\n")) + 2
				break
			}
		}

		if bodyStartIndex < n {
			body := buffer[bodyStartIndex:n]
			err := os.WriteFile(filePath, body, 0644)
			if err != nil {
				fmt.Println("Error writing file: ", err.Error())
				response = "HTTP/1.1 500 Internal Server Error\r\n\r\n"
			} else {
				response = "HTTP/1.1 201 Created\r\n\r\n"
			}
		} else {
			response = "HTTP/1.1 400 Bad Request\r\n\r\n"
		}
	} else {
		switch subRoute[1] {
		case "":
			response = "HTTP/1.1 200 OK\r\n\r\n"
		case "user-agent":
			ua := strings.Split(lines[2], " ")[1]
			response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(ua), ua)
		case "echo":
			body := strings.Join(subRoute[2:], "/")
			response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(body), body)
		case "files":
			file_name := strings.Join(subRoute[2:], "/")
			filePath := filepath.Join(directory, file_name)
			response = readFileContents(filePath)
		default:
			response = "HTTP/1.1 404 Not Found\r\n\r\n"
		}
	}

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error sending msg: ", err.Error())
		//os.Exit(1)
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	var directory string
	flag.StringVar(&directory, "directory", "", "the directory to serve files from")
	flag.Parse()
	
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
		}
		go handleConnection(conn, directory)
	}
}
