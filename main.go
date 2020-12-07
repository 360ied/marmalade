package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"marmalade/config"
	"marmalade/packets/inbound"
	"marmalade/packets/outbound"
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

	_ /* protocol version */, username, _ /* verification key */, readPlayerIdentificationErr := inbound.ReadPlayerIdentification(reader)
	if readPlayerIdentificationErr != nil {
		log.Printf("ERROR: Error reading player identification packet from %v, error: %v", conn.RemoteAddr().String(), readPlayerIdentificationErr)
		return
	}
	log.Printf("INFO: Recieved a player identification packet from %v, they say their username is `%v`", conn.RemoteAddr().String(), username)

	writer := outbound.NewAFCBW(conn, config.BufferFlushInterval)
	defer writer.Close()

	sendServerIdentificationErr := writer.SendServerIdentification(config.ServerName, config.ServerMOTD, false)
	if sendServerIdentificationErr != nil {
		log.Printf("ERROR: Error sending server identification packet to %v, error: %v", conn.RemoteAddr().String(), sendServerIdentificationErr)
		return
	}
}
