package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	pb "spacegineers_context/protobufs/room_proto"

	"google.golang.org/grpc"
)

type Supply struct {
	supplyType        string
	supplyAvalability bool
}

// GRPC Server implementaion
type Service struct {
	pb.UnimplementedRoomGRPCServer
	supply chan *Supply
}

func (s *Service) Electricity(ctx context.Context, source *pb.Source) (empty *pb.Empty, err error) {
	fmt.Printf("Received electricity: %v \n", source.On)
	s.supply <- &Supply{
		supplyType:        "elec",
		supplyAvalability: source.On,
	}
	return empty, nil
}

func (s *Service) Air(ctx context.Context, source *pb.Source) (empty *pb.Empty, err error) {
	fmt.Printf("Received air: %v \n", source.On)
	s.supply <- &Supply{
		supplyType:        "air",
		supplyAvalability: source.On,
	}
	return empty, nil
}

func (s *Service) ClientRegister(stream pb.RoomGRPC_ClientRegisterServer) error {
	fmt.Println("A client registered")
	recError := make(chan error)

	go func() {
		for {
			rec, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("connections lost with the room client")
				recError <- errors.New("connections lost with the room client")
				return
			}
			if err != nil {
				fmt.Println("Failed to receive message from the room client")
				recError <- err
				return
			}

			fmt.Printf("Received damage from client: %d \n", rec.Damage)
		}
	}()

	for {
		select {
		case err := <-recError:
			return err
		case supple := <-s.supply:
			sendSupply := &pb.Supply{
				Type: supple.supplyType,
				On:   supple.supplyAvalability,
			}
			if err := stream.Send(sendSupply); err != nil {
				fmt.Println("Failed to send supply message to the room client")
				recError <- err
			}
		}
	}
}

func main() {
	address := "localhost:53000"

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen to the address %s: %v", address, err)
	}

	grpcServer := grpc.NewServer()

	service := &Service{supply: make(chan *Supply)}

	pb.RegisterRoomGRPCServer(grpcServer, service)
	log.Println("Room Server: Started")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
