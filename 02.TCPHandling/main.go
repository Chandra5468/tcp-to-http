package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {

	// Create or open tcp.txt for writing output (overwrites existing file)

	outFile, err := os.Create("tcp.txt")

	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	// Tee output to both stdout and file

	multiWriter := io.MultiWriter(os.Stdout, outFile)
	log.SetOutput(multiWriter) // send log messages to both console and file

	// Use net.Listen to setup a TCP listener on PORT :42069
	listerner, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}

	defer listerner.Close()
	log.Println("TCP Server listening on :42069")

	for {
		conn, err := listerner.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		log.Println("A connection has been accepted")

		go func(c net.Conn) { // Conn is a generic stream-oriented network connection.
			fmt.Fprintf(multiWriter, "Handling connection from %v\n", conn.LocalAddr())
			defer func() {
				c.Close()
				fmt.Fprintf(multiWriter, "Connection closed for %v\n", c.RemoteAddr())
			}()

			lines := getLinesChannel(c)
			for line := range lines {
				fmt.Fprintln(multiWriter, line)
			}
		}(conn)
	}
}

func getLinesChannel(conn net.Conn) <-chan string {
	lines := make(chan string)

	go func() {
		defer close(lines)
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
	}()

	return lines
}
