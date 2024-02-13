/*
This file contains information about a Room

A Room is a grpc client
*/

package main

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

type Room struct {
	controllerAddress     string
	componentId           string
	componentName         string
	ctxIdentifier         string
	resourcesAvailability map[string]bool
}

func NewRoom(addr string, ctxId string, identifier string, name string) *Room {
	return &Room{
		controllerAddress:     addr,
		componentId:           identifier,
		componentName:         name,
		ctxIdentifier:         ctxId,
		resourcesAvailability: make(map[string]bool),
	}
}

func (c *Room) StartServer() {
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
		fmt.Printf("%s : receives {%v} from the Room Controller \n", c.componentName, recSupply)

		// if receive availability is the same as the stored availability
		if !c.resourcesAvailability[recSupply.ResourceType] {
			c.resourcesAvailability[recSupply.ResourceType] = true
			if err := stream.Send(recSupply); err != nil {
				log.Fatalf("%s : failed to send message to the Room Controller, %v", c.componentName, err)
			}
			fmt.Printf("%s : send {%v} to the Room Controller \n", c.componentName, recSupply)
		}
	}
}

func closeConnection(conn *grpc.ClientConn, name string) {
	if err := conn.Close(); err != nil {
		log.Fatalf("%s : failed to end communication with the Room Controller, %v", name, err)
	}
}
