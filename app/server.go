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
	// request := string(buffer)
	request := string(buffer[:n])
	lines := strings.Split(request, "\r\n")
	start_line := lines[0]
	// headers := lines[1]
	start_line_parts := strings.Fields(start_line)
	rMethod, rPath, rProtocol := start_line_parts[0], start_line_parts[1], start_line_parts[2]
	fmt.Printf("method=%s, path=%s, protocol=%s\n", rMethod, rPath, rProtocol)
	fmt.Printf("lines: %s\n", lines)
	fmt.Printf("lines0: %s\n", lines[0])
	fmt.Printf("lines1: %s\n", lines[1])
	fmt.Printf("lines2: %s\n", lines[2])

	subRoute := strings.Split(rPath, "/")
	// ua := strings.Split(lines[2], " ")[1]

	var response string

	if rMethod == "POST" {
        fileName := strings.TrimPrefix(rPath, "/files/")
        filePath := filepath.Join(directory, fileName)

        // Find the index of the empty line
        var headerLength int
        for i, line := range lines {
            if line == "" {
                headerLength = i
                break
            }
        }

        // Assuming headers and body fit in the buffer, extract the body
        // body := strings.Join(lines[headerLength+1:], "\r\n")
		body := strings.Join(lines[headerLength+1:], "\r\n")
		actualBody := []byte(body)
        err := os.WriteFile(filePath, []byte(actualBody), 0644)
        if err != nil {
            fmt.Println("Error writing file: ", err.Error())
            response = "HTTP/1.1 500 Internal Server Error\r\n\r\n"
        } else {
            response = "HTTP/1.1 201 Created\r\n\r\n"
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
