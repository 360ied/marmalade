package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"marmalade/commands"
	"marmalade/config"
	"marmalade/packets/inbound"
	"marmalade/packets/outbound"
	"marmalade/world"
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

	p := &world.Player{
		Username: username,
		OP:       false,
		Writer:   writer,
	}
	if !world.AddPlayer(p) {
		log.Printf("ERROR: Max players reached!")
		return
	}
	defer world.RemovePlayer(p.ID)
	log.Printf("INFO: Assigned `%v` player id %v", username, p.ID)

	if err := world.SendWorld(writer); err != nil {
		log.Printf("ERROR: Failed to send world: %v", err)
		return
	}

	if err := writer.SendSpawnPlayer(
		255,
		username,
		config.SpawnPointX,
		config.SpawnPointY,
		config.SpawnPointZ,
		config.SpawnPointYaw,
		config.SpawnPointPitch); err != nil {
		log.Printf("ERROR: Failed to send spawn player: %v", err)
		return
	}

	world.SpawnOtherPlayers(p)

	for {
		b, bErr := reader.ReadByte()
		if bErr != nil {
			log.Printf("ERROR: Failed to read packet id: %v", bErr)
			return
		}
		if err := reader.UnreadByte(); err != nil {
			log.Printf("ERROR: Failed to unread packet id: %v", err)
			return
		}
		switch b {
		case 0x05: // set block
			x, y, z, mode, blockType, packetErr := inbound.ReadSetBlock(reader)
			if packetErr != nil {
				log.Printf("ERROR: Failed to read set block packet: %v", packetErr)
				return
			}
			world.HandleSetBlock(x, y, z, mode, blockType)
		case 0x08: // position and orientation
			x, y, z, yaw, pitch, packetErr := inbound.ReadPositionAndOrientation(reader)
			if packetErr != nil {
				log.Printf("ERROR: Failed to read position and orientation packet: %v", packetErr)
				return
			}
			world.HandlePositionAndOrientation(p, x, y, z, yaw, pitch)
		case 0x0d: // message
			message, packetErr := inbound.ReadMessage(reader)
			if packetErr != nil {
				log.Printf("ERROR: Failed to read chat message: %v", packetErr)
				return
			}
			if strings.HasPrefix(message, config.CommandPrefix) {
				commands.HandleCommand(p, message[len(config.CommandPrefix):])
			} else {
				world.BroadcastMessage(fmt.Sprintf("<%v> %v", p.Username, message))
			}
		}
	}
}
