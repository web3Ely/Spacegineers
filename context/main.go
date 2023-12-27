package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "spacegineers_context/protobufs"
)

const (
	port = 8080
)

type SpaceAndTime struct {
	pb.UnimplementedServicesServer
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("Failed to listen to the port %d: %v", port, err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterServicesServer(grpcServer, &SpaceAndTime{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
