package dvs

import (
	"bytes"
	"fmt"
	"log"
	"net"
)

// Master is a master client struct, responsible for storing
// config data for DVS connection. On Connect it spawns
// the number of Clients equal to the number of SourceIds
// provided, and connects them to DVS server.
type Master struct {
	SourceIds       []int
	DestinationId   int
	MopPpid         int
	Address         string
	CommandsChannel chan Command
}

func (master Master) Connect() {
	var sourceId int
	for _, sourceId = range master.SourceIds {
		go Client{Master: master, SourceId: sourceId}.Connect()
	}
}

type Client struct {
	Master
	SourceId      int
	CommandsQueue map[int]Command
	// state ClientState (connected, initialized, ready, closed)
	// IF closed then try to reconnect, If connected then try to initialize,
	// if initialized then try to by ready?
}

func (client Client) Connect() {
	client.CommandsQueue = make(map[int]Command)
	log.Printf("Client[%d] spawned", client.SourceId)

reconnect:
	connection, err := net.Dial("tcp", client.Address)
	if err != nil {
		log.Fatalf("Client[%d] CONNECTION ERROR %s", client.SourceId, err)
	}
	log.Printf("Client[%d] connected", client.SourceId)

	err = client.Initialize(&connection)
	if err != nil {
		log.Fatalf("Client[%d] connection initialization error: %s", client.SourceId, err)
		// TODO(m): what now? retry? reject?
		goto reconnect
		// What if it crashes at reading?
	}

	for {
		command := <-client.CommandsChannel
		log.Printf("Client[%d] recieved command", client.SourceId)
		transactionId, request := client.PrepareRequest(command)
		log.Printf("Client[%d] prepared request (%dbytes) %v", client.SourceId, len(request), request)
		client.CommandsQueue[transactionId] = command
		client.SendRequest(&connection, request)
	}
}

func (client Client) Initialize(connection *net.Conn) error {
	//(*connection).SetDeadline(30)
	client.SendRequest(connection, Message1)
	err := client.ExpectMessage2Success(connection)
	if err != nil {
		return err
	}
	return client.ExpectMessage3CallAccepted(connection)
}

func (client Client) ExpectMessage2Success(connection *net.Conn) error {
	if size, err := client.ReadSize(connection); err == nil {
		message, err := client.ReadMessage(connection, size)
		if bytes.Equal(message, Message2Success) {
			return err
		}
	}
	return nil
}

func (client Client) ExpectMessage3CallAccepted(connection *net.Conn) error {
	if size, err := client.ReadSize(connection); err == nil {
		message, err := client.ReadMessage(connection, size)
		if bytes.Equal(message, Message3CallAccepted) {
			return err
		}
	}
	return nil
}

func (client Client) ReadSize(connection *net.Conn) (int, error) {
	buf := make([]byte, 2)
	// Read exactly 2 bytes
	_, err := (*connection).Read(buf)
	if err != nil {
		log.Fatalf("Client[%d] read error: %s", client.SourceId, err)
		return 0, err
	}
	size, err := DecodeHex(buf)
	if err != nil {
		log.Fatalf("Client[%d] size decode error: %s", client.SourceId, err)
		return 0, err
	}
	return int(size), nil
}

func (client Client) ReadMessage(connection *net.Conn, size int) ([]byte, error) {
	buf := make([]byte, size)
	rsize, err := (*connection).Read(buf)
	if (err != nil) || (rsize != size) {
		log.Fatalf("Client[%d] read error(size=%d expected=%d): %s", client.SourceId, rsize, size, err)
		return []byte{}, err
	}
	return buf, nil
}

func (client *Client) PrepareRequest(command Command) (transactionId int, request []byte) {
	header := RootHeader(transactionId, command.Type, client.SourceId, client.DestinationId, client.MopPpid)
	request = DeviceIO(fmt.Sprint(header, command.Body))
	transactionId = 1
	return
}

func (client *Client) SendRequest(connection *net.Conn, request []byte) {
	length, err := (*connection).Write(request)
	if err != nil {
		log.Fatalf("Client[%d] WRITE ERROR %s", client.SourceId, err)
	}
	log.Printf("Client[%d] bytes send %d", client.SourceId, length)

	if length != len(request) {
		log.Fatalf("Client[%d] write failed. %d out of %d sent.", client.SourceId, length, len(request))
	}

	log.Printf("Client[%d] sent command", client.SourceId)
}
