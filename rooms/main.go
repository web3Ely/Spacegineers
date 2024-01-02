package main

import (
	"context"
	"fmt"
	"io"
	"log"
	pb "spacegineers_context/protobufs"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	const (
		port = 8080
	)
	log.Println("Room started")

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Bridge could not connect to context: %v", err)
	}
	defer conn.Close()

	// Create a client
	client := pb.NewSpaceContextClient(conn)

	// Client sends a stream of requests
	stream, err := client.RoomRegister(context.Background(), &pb.Nothing{})
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	// Client can receive messages from the server
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			log.Println("Stream is over")
			break
		}
		if err != nil {
			log.Fatalf("Error receiving message: %v", err)
		}

		// Process the received response
		fmt.Println("Received:", response)
	}
}
