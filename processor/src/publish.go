package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
)

var (
	clientConn net.Conn
	connOnce   sync.Once
	connErr    error
)

// fazer a reconexão automatica
// caso a conexão caia tenha um fallback
func getClient() (net.Conn, error) {
	connOnce.Do(func() {
		clientConn, connErr = net.Dial("tcp", "messaging:4555")
	})
	return clientConn, connErr
}

func Publish(correlationId string, amount float64) error {
	client, err := getClient()
	if err != nil {
		fmt.Println(err)
		return err
	}

	msg, err := marshall(correlationId, amount)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = client.Write(msg)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func marshall(correlationId string, amount float64) ([]byte, error) {
	buffer := new(bytes.Buffer)

	buffer.WriteString(correlationId)

	err := binary.Write(buffer, &binary.BigEndian, amount)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return buffer.Bytes(), nil
}
