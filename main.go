package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"marmalade/config"
)

func main() {
	// Create new listener
	listener, listenerErr := net.Listen("tcp", config.Address)
	if listenerErr != nil {
		panic(fmt.Sprintf("FATAL: Error starting TCP listener: %v", listenerErr))
	}
	// Close listener on exit
	defer func() {
		if err := listener.Close(); err != nil {
			log.Panicf("FATAL: Error closing TCP listener: %v", err)
		} else {
			log.Println("INFO: Successfully closed TCP listener.")
		}
	}()
	// Accept connections
	for {
		connection, connectionErr := listener.Accept()
		if connectionErr != nil {
			log.Printf("ERROR: Failed to accept TCP connection: %v\n", connectionErr)
			continue
		}
		go handleConnection(connection)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("ERROR: Failed to close TCP connection with %v", conn.RemoteAddr().String())
		} else {
			log.Printf("INFO: Sucessfully closed TCP connection with %v", conn.RemoteAddr().String())
		}
	}()
	log.Printf("INFO: %v has established a connection.", conn.RemoteAddr().String())

	reader := bufio.NewReader(conn)

}
