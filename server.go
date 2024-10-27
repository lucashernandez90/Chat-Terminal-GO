package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type client struct {
	name string
	ch   chan<- string // canal de mensagem
}

var (
	entering   = make(chan client)
	leaving    = make(chan client)
	messages   = make(chan string)
	private    = make(chan privateMessage)
	updateNick = make(chan updateNickname)
)

type privateMessage struct {
	from    string
	to      string
	message string
}

type updateNickname struct {
	oldName string
	newName string
	cli     client
	result  chan bool // Canal para confirmar o resultado da troca
}

func broadcaster() {
	clients := make(map[client]bool)
	clientNames := make(map[string]client)
	for {
		select {
		case msg := <-messages:
			// Envia a mensagem para todos os clientes
			for cli := range clients {
				cli.ch <- msg
			}
		case pmsg := <-private:
			// Mensagem privada
			if recipient, ok := clientNames[pmsg.to]; ok {
				recipient.ch <- "[Privado] " + pmsg.from + ": " + pmsg.message
			} else {
				if sender, ok := clientNames[pmsg.from]; ok {
					sender.ch <- "Usuário @" + pmsg.to + " não encontrado."
				}
			}
		case update := <-updateNick:
			// Atualizar apelido
			if _, exists := clientNames[update.newName]; exists {
				// Apelido já em uso
				update.cli.ch <- "O apelido @" + update.newName + " já está em uso."
				update.result <- false
			} else {
				// Atualiza o mapa de clientes
				delete(clientNames, update.oldName)
				update.cli.name = update.newName
				clientNames[update.newName] = update.cli

				// Notifica todos os usuários sobre a mudança
				messages <- "Usuário @" + update.oldName + " agora é @" + update.newName
				update.result <- true
			}
		case cli := <-entering:
			// Adiciona um novo cliente
			clients[cli] = true
			clientNames[cli.name] = cli
		case cli := <-leaving:
			// Remove um cliente
			delete(clients, cli)
			delete(clientNames, cli.name)
			close(cli.ch)
		}
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	input := bufio.NewScanner(conn)

	var apelido string
	ch <- "Por favor, escolha um apelido:"
	if input.Scan() {
		apelido = input.Text()
	}

	cli := client{name: apelido, ch: ch}
	messages <- "Usuário @" + apelido + " acabou de entrar"
	entering <- cli

	for input.Scan() {
		msg := input.Text()
		if strings.HasPrefix(msg, "\\changenick ") {
			newNickname := msg[len("\\changenick "):]
			messages <- "Usuário @" + apelido + " mudou seu nome para @" + newNickname
			cli.name = newNickname // Atualiza o nome no cliente
			apelido = newNickname  // Atualiza a variável de apelido
		} else if strings.HasPrefix(msg, "\\msg ") {
			parts := strings.SplitN(msg[5:], " ", 2)
			if len(parts) == 2 {
				private <- privateMessage{from: apelido, to: strings.TrimPrefix(parts[0], "@"), message: parts[1]}
			}
		} else {
			messages <- "@" + apelido + " disse: " + msg
		}
	}

	leaving <- cli
	messages <- "Usuário @" + apelido + " saiu"
	conn.Close()
}

func main() {
	fmt.Println("Iniciando servidor...")
	listener, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
