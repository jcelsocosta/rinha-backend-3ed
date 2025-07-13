package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func RunMessagingServer() {
	fmt.Println("Messaging on :4555")
	listener, err := net.Listen("tcp", ":4555")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		message, err := unmarshall(reader)
		if err != nil {
			continue
		}

		select {
		case channelMsg <- *message:
			// enviado com sucesso
		default:
			log.Println("Erro: canal cheio, mensagem descartada:", message.CorrelationId)
		}
	}
}
