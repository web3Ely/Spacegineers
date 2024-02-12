/*
This file contains information about a component type

A component type represents a component client in grpc paradigm which can then
be initialized by other actural compartment of a ship object like room and appliance.
*/

package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	pb "spacegineers_context/protobufs/controller_proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Component struct {
	controllerAddress string
	componentId       string
	componentName     string
	ctxIdentifier     string
}

func NewComponent(addr string, identifier string, name string, ctxId string) *Component {
	return &Component{
		controllerAddress: addr,
		componentId:       identifier,
		componentName:     name,
		ctxIdentifier:     ctxId,
	}
}

func (c *Component) StartServer() {
	connection, err := grpc.Dial(c.controllerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer closeConnection(connection, c.componentName)

	if err != nil {
		log.Fatalf("%s : failed to connect to the Room Controller, %v", c.componentName, err)
	}
	client := pb.NewControllerGRPCClient(connection)
	fmt.Printf("%s : successfully connected to the Room Controller \n", c.componentName)

	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(c.ctxIdentifier, c.componentId))
	stream, err := client.Subscribe(ctx)
	if err != nil {
		log.Fatalf("%s : failed to subscribe to the Room Controller, %v", c.componentName, err)
	}
	fmt.Printf("%s : successfully registered to the Room Controller \n", c.componentName)

	for {
		recSupply, err := stream.Recv()
		if err == io.EOF {
			// when server return nil
			log.Fatalf("%s : connection dropped by the Room Controller, %v", c.componentName, err)
		}
		if err != nil {
			log.Fatalf("%s : failed to receive message from the Room Controller, %v", c.componentName, err)
		}
		fmt.Printf("%s : received %v from the Room Controller \n", c.componentName, recSupply)

		if err := stream.Send(recSupply); err != nil {
			log.Fatalf("%s : failed to send message to the Room Controller, %v", c.componentName, err)
		}
	}
}

func closeConnection(conn *grpc.ClientConn, name string) {
	if err := conn.Close(); err != nil {
		log.Fatalf("%s : failed to end communication with the Room Controller, %v", name, err)
	}
}
