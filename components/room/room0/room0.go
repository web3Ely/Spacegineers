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

// It establish server side streaming with context api
// It receives damage from context api and send it to RoomToServer
func main() {
	controller := "localhost:53000"
	name := "Room0"
	ctxIdentifier := "compcode"
	componentNum := "R0"

	connection, err := grpc.Dial(controller, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("%s failed to connect to the Room Controller, %v", name, err)
	}
	defer connection.Close()

	client := pb.NewControllerGRPCClient(connection)
	fmt.Printf("%s successfully connected to the Room Controller \n", name)

	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(ctxIdentifier, componentNum))
	stream, err := client.Subscribe(ctx)
	if err != nil {
		log.Fatalf("%s failed to subscribe to the Room Controller, %v", name, err)
	}
	fmt.Printf("%s successfully registered to the Room Controller \n", name)

	for {
		rec, err := stream.Recv()
		if err == io.EOF {
			// when server return nil
			log.Fatalf("%s connection dropped by the Room Controller, %v", name, err)
		}
		if err != nil {
			log.Fatalf("%s failed to communicate with the Room Controller, %v", name, err)
		}
		fmt.Printf("%s received %v from the Room Controller \n", name, rec)

	}
}
