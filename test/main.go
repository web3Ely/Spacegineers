package main

import (
	"context"
	"fmt"
	"log"

	pb "spacegineers_context/protobufs"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	port = 8080
)

func main() {
	log.Println("Starting a test unary grpc call to the context api")
	connection, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Test grpc call failed: %v", err)
	}
	client := pb.NewSpaceContextClient(connection)
	log.Println("Test client creation successful")

	for i := 0; i < 3; i++ {
		log.Println("----Sending a grpc call")
		client.Hit(context.Background(), &pb.Damage{Damage: 1000})
		log.Println("call ends")
	}

	log.Println("Finish the test call")
}
