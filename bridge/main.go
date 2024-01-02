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

const (
	port = 8080
)

func main() {
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Bridge could not connect to context: %v", err)
	}
	defer conn.Close()

	client := pb.NewServicesClient(conn)

	res, err := client.ContextService(context.Background(), &pb.DamageRequest{Damage: 1000})
	if err != nil {
		log.Fatalf("Could not send damage to context %v", err)
	}
	log.Printf("Response is: %s", res)

	stream, err := client.StreamData(context.Background(), &pb.NullValue{})
	if err != nil {
		log.Fatalf("Error opening stream: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Println("It is end of file")
			break
		}
		if err != nil {
			log.Fatalf("Error: %v.StreamData(_) = _, %v", client, err)
		}

		log.Printf("Received integer: %d\n", resp.GetResult())
	}

	log.Println("Stream closed.")
}
