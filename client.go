package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Conectado ao servidor!")

	done := make(chan struct{})

	go func() {
		io.Copy(os.Stdout, conn)
		done <- struct{}{}
	}()

	inputReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, _ := inputReader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "\\exit" {
			fmt.Fprintln(conn, input)
			conn.Close()
			break
		} else {
			fmt.Fprintln(conn, input)
		}
	}

	<-done
}
