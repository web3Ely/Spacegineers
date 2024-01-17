package main

import (
	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Error channel that let the main know something is wrong in communication with different services
	// damage channel carries damage info received by context api to server call
	// closeChan channel allows bidirectional communication between client and server to close its stream when Error channel receives
	roomError := make(chan error)
	damage := make(chan int32)
	roomClose := make(chan bool)
	var wg sync.WaitGroup
	defer close(roomError)
	defer close(damage)
	defer close(roomClose)

	// Set up all the addresses
	adrs := &ConnectionAdrs{
		serverAdr:           "localhost:53000",
		contextAdr:          "localhost:8080",
		downstreamRoomsAdrs: []string{},
	}

	// Creat connections to all the addresses
	connections := adrs.Connect()
	defer connections.CloseAll()

	// Set up all the connections
	destinies := &Destinies{
		server:          connections.connections[0],
		context:         connections.connections[1],
		downstreamRooms: connections.connections[2:],
	}

	// Server should start first
	wg.Add(2)
	go ConnectServer(destinies, roomError, damage, roomClose, &wg)

	go ConnectContex(destinies, roomError, damage, roomClose, &wg)

	// Close all connection and clien-server stream and log an error
	err := <-roomError
	roomClose <- true
	wg.Wait()
	defer log.Printf("Error: %v", err)
}

type Destinies struct {
	server          *grpc.ClientConn
	context         *grpc.ClientConn
	downstreamRooms []*grpc.ClientConn
}

type Connections struct {
	connections []*grpc.ClientConn
}

func (c *Connections) GetConnection(adr string) {
	conn, err := grpc.Dial(adr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.CloseAll()
		log.Fatalf("Could not connect to %s, error: %v", adr, err)
	}
	c.connections = append(c.connections, conn)
}

func (c *Connections) CloseAll() {
	for _, conn := range c.connections {
		conn.Close()
	}
}

type ConnectionAdrs struct {
	serverAdr           string
	contextAdr          string
	downstreamRoomsAdrs []string
}

func (adrs *ConnectionAdrs) Connect() *Connections {
	connections := &Connections{
		connections: []*grpc.ClientConn{},
	}
	connections.GetConnection(adrs.serverAdr)
	connections.GetConnection(adrs.contextAdr)
	for _, adr := range adrs.downstreamRoomsAdrs {
		connections.GetConnection(adr)
	}
	return connections
}
