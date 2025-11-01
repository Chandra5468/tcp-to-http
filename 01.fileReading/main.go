package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer close(lines)
		defer f.Close()

		buffer := make([]byte, 20)

		currentLine := ""

		for {
			bytesRead, err := f.Read(buffer)

			if err != nil {
				if err == io.EOF {
					if len(currentLine) > 0 {
						lines <- currentLine
					}
					break
				}
				log.Println("error reading file: ", err)
				break
			}

			currentLine += string(buffer[:bytesRead])
			parts := strings.Split(currentLine, "\n")

			for _, part := range parts[:len(parts)-1] {
				lines <- part
				// time.Sleep(time.Second)
			}

			currentLine = parts[len(parts)-1]
		}
	}()

	return lines
}

func main() {
	f, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatal(err)
	}

	for line := range getLinesChannel(f) {
		fmt.Println("Line:", line)
	}

	// TCP Learning Below
	test()
}
