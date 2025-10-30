package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	f, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	buffer := make([]byte, 8)
	currentLine := ""
	for {
		bytesRead, err := f.Read(buffer)
		if err != nil {
			if err == io.EOF {
				// Handle any remaining data in currentLine at EOF
				if len(currentLine) > 0 {
					fmt.Println("Line:", currentLine)
				}
				break
			}
			fmt.Println("error reading file: ", err)
			return
		}

		// Concatenate, split by newline
		currentLine += string(buffer[:bytesRead])
		parts := strings.Split(currentLine, "\n")
		// Print all lines except last
		for _, part := range parts[:len(parts)-1] {
			fmt.Println("Line:", part)
			time.Sleep(time.Second * 1)
		}
		// Last part could be incomplete, keep for next round
		currentLine = parts[len(parts)-1]
	}
}
