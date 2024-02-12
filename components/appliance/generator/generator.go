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

// It establish bidirectional stream with room server
// It also send explosion signal to context api
func main() {
	controller := "localhost:53000"
	componentName := "Generator"
	componentId := "G0"
	ctxIdentifier := "compcode"

	connection, err := grpc.Dial(controller, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("%s failed to connect to the Room Controller, %v", componentName, err)
	}
	client := pb.NewControllerGRPCClient(connection)
	fmt.Printf("%s successfully connected to the Room Controller \n", componentName)

	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(ctxIdentifier, componentId))
	stream, err := client.Subscribe(ctx)
	if err != nil {
		log.Fatalf("%s failed to subscribe to the Room Controller, %v", componentName, err)
	}
	fmt.Printf("%s successfully registered to the Room Controller \n", componentName)

	supply := &pb.Supply{
		SenderId:            componentId,
		ResourceType:        "elec",
		ResouceAvailability: true,
	}
	if err := stream.Send(supply); err != nil {
		log.Fatalf("%s initiate message to the Room Controller failed, %v", componentName, err)
	}

	for {
		rec, err := stream.Recv()
		if err == io.EOF {
			// when server return nil
			log.Fatalf("%s connection dropped by the Room Controller, %v", componentName, err)
		}
		if err != nil {
			log.Fatalf("%s failed to communicate with the Room Controller, %v", componentName, err)
		}
		fmt.Printf("%s received %v from the Room Controller \n", componentName, rec)
	}
}
