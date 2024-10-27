package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func main() {

	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		fmt.Println("Erro ao se conectar ao servidor:", err)
		return
	}
	defer conn.Close()

	botName := "Bot@ReverseBot"
	fmt.Fprintln(conn, botName)

	reader := bufio.NewReader(conn)
	for {

		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		if strings.Contains(message, "disse:") && !strings.Contains(message, botName) {
			parts := strings.Split(message, ": ")
			if len(parts) > 1 {
				receivedMsg := parts[1]
				reversed := reverse(receivedMsg)
				fmt.Println("Recebi:", receivedMsg)
				fmt.Println("Resposta:", reversed)
				fmt.Fprintln(conn, reversed)
			}
		}
	}
}
