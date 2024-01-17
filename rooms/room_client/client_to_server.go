package main

import (
	ctx "context"
	"fmt"
	"io"
	pbServer "spacegineers_context/protobufs/room_proto"

	// "spacegineers_context/rooms/room_utility"
	"sync"
)

func ConnectServer(des *Destinies, roomError chan<- error, damage <-chan int32, roomClose <-chan bool, wg *sync.WaitGroup) {
	recError := make(chan error, 1)
	defer close(recError)

	client := pbServer.NewRoomGRPCClient(des.server)
	// downstreamClients := make([]pbServer.RoomGRPCClient, len(des.downstreamRooms))

	stream, err := client.ClientRegister(ctx.Background())
	if err != nil {
		fmt.Println("Client to Server: Error registering to the room server")
		roomError <- err
		return
	}
	defer stream.CloseSend()

	// myStream := room_utility.MyStream{
	// 	Width: 0,
	// }

	// room_utility.Receive[*pbServer.RoomGRPC_ClientRegisterClient](&stream)

	fmt.Println("Client to Server: Finish registration to the room server")

	// Receive message from the server and send them to downstream rooms
	go func() {
		for {
			rec, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Client to Server: Connection lost with the room server")
				recError <- err
				return
			}
			if err != nil {
				if _, ok := <-recError; ok {
					fmt.Println("Client to Server: Error receiving message from the room server")
					recError <- err
				}

				return
			}
			fmt.Println(rec.On)
			// for _, downstreamClient := range downstreamClients {
			// 	switch rec.Type {
			// 	case "elec":
			// 		if _, err := downstreamClient.Electricity(ctx.Background(), &pbServer.Source{On: rec.On}); err != nil {
			// 			fmt.Println("Client to Server: Failed to send message to the downstream room")
			// 		}
			// 	case "air":
			// 		if _, err := downstreamClient.Air(ctx.Background(), &pbServer.Source{On: rec.On}); err != nil {
			// 			fmt.Println("Client to Server: Failed to send message to the downstream room")
			// 		}
			// 	}
			// }
		}
	}()

	defer wg.Done()

	for {
		select {
		case err := <-recError:
			roomError <- err
			return
		case <-roomClose:
			return
		case d := <-damage:
			fmt.Println("Client to Server: Sending damge to the room server")
			if err := stream.Send(&pbServer.Damage{Damage: d}); err != nil {
				fmt.Println("Client to Server: Failed to send damage message to the room server")
				recError <- err
			}
		}
	}
}
