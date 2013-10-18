package dvs

import (
	"fmt"
	"log"
	"net"
)

// Master is a master client struct, responsible for storing
// config data for DVS connection. On Connect it spawns
// the number of Clients equal to the number of SourceIds
// provided, and connects them to DVS server.
type Master struct {
	SourceIds []int
	DestinationId int
	MopPpid int
	Address string
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
	SourceId int
	CommandsQueue map[int]Command
}

func (client Client) Connect() {
	client.CommandsQueue = make(map[int]Command)
	log.Printf("Client[%d] spawned", client.SourceId)

	connection, err := net.Dial("tcp", client.Address)
	if err != nil {
		log.Fatalf("Client[%d] CONNECTION ERROR %s", client.SourceId, err)
	}
	log.Printf("Client[%d] connected", client.SourceId)

	for {
		command := <-client.CommandsChannel
		log.Printf("Client[%d] recieved command", client.SourceId)
		transactionId, request := client.PrepareRequest(command)
		client.CommandsQueue[transactionId] = command
		client.SendRequest(&connection, request)
	}
}

func (client *Client) PrepareRequest(command Command) (transactionId int, request []byte){
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

	if length != len(request) {
		log.Fatalf("Client[%d] write failed. %d out of %d sent.", client.SourceId, length, len(request))
	}

	log.Printf("Client[%d] sent command", client.SourceId)
}
