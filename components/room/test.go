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

// It establish server side streaming with context api
// It receives damage from context api and send it to RoomToServer
func main() {
	controller := "localhost:53000"
	ctxIdentifier := "itemcode"
	roomNum := "1"

	connection, err := grpc.Dial(controller, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("test failed to connect to the Room Controller, %v", err)
	}
	client := pb.NewRoomGRPCClient(connection)
	fmt.Println("test successfully connected to the Room Controller")

	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(ctxIdentifier, roomNum))
	stream, err := client.RoomSubscribe(ctx)
	if err != nil {
		log.Fatalf("test failed to subscribe to the Room Controller, %v", err)
	}
	fmt.Println("test successfully registered to the Room Controller")

	stream.Send(&pb.Supply{})
	for {
		rec, err := stream.Recv()
		if err == io.EOF {
			// when server return nil
			log.Fatalf("test connection dropped by the Room Controller, %v", err)
		}
		if err != nil {
			log.Fatalf("test failed to communicate with the Room Controller, %v", err)
		}
		fmt.Printf("Room1 received %v from the Room Controller \n", rec)

	}
}
