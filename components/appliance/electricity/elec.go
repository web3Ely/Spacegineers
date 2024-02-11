package main

import (
	"context"
	"fmt"
	"io"
	"log"
	pb "spacegineers_context/protobufs/room_proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// It establish bidirectional stream with room server
// It also send explosion signal to context api
func main() {
	controller := "localhost:53000"
	name := "Generator"
	ctxIdentifier := "compcode"
	componentNum := "2"

	connection, err := grpc.Dial(controller, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("%s failed to connect to the Room Controller, %v", name, err)
	}
	client := pb.NewRoomGRPCClient(connection)
	fmt.Printf("%s successfully connected to the Room Controller \n", name)

	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(ctxIdentifier, componentNum))
	stream, err := client.RoomSubscribe(ctx)
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
