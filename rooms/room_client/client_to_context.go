package main

import (
	ctx "context"
	"fmt"
	"io"
	pbContex "spacegineers_context/protobufs/context_proto"

	// "spacegineers_context/rooms/room_utility"
	"sync"
)

func ConnectContex(des *Destinies, roomError chan<- error, damage chan<- int32, roomClose <-chan bool, wg *sync.WaitGroup) {
	recError := make(chan error)
	damageReceived := make(chan int32)
	defer close(recError)
	defer close(damageReceived)

	client := pbContex.NewSpaceContextClient(des.context)

	stream, err := client.RoomRegister(ctx.Background(), &pbContex.Empty{})
	if err != nil {
		fmt.Println("Client to Context: Error registering to the context api")
		roomError <- err
		return
	}
	defer stream.CloseSend()

	// room_utility.Receive[*pbContex.SpaceContext_RoomRegisterClient](&stream)

	fmt.Println("Client to Context: Finish registration to the context api")

	go func() {
		for {
			rec, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Client to Context: Connection lost with the context api")
				recError <- err
				return
			}
			if err != nil {
				fmt.Println("Client to Context: Error receiving message from the context api")
				recError <- err
				return
			}
			fmt.Println("Client to Context: Received Damge form the context api")
			damageReceived <- rec.Damage
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
		case dmg := <-damageReceived:
			damage <- dmg
		}
	}
}
