package controllerhelper

import (
	"fmt"
	"io"
	"log"
	"net"
	"slices"
	pb "spacegineers_context/protobufs/room_proto"

	// shipSetting "spacegineers_context/settings"

	"google.golang.org/grpc"
)

// Incoming and outgoing data
// defined in protobuf
type Supply struct {
	senderId            int32
	resourceType        string
	resouceAvailability bool
}

// Controller is a GRPC Server type
type Controller struct {
	pb.UnimplementedRoomGRPCServer
	Address             string
	ComponentIdentifier string
	ConnectionTable     map[int32][]int32
	pathways            map[int32]chan<- *Supply
}

// Controller constructor
func NewController() *Controller {
	return &Controller{}
}

// This method establish server-side streamming between each room and Controller
func (c *Controller) RoomSubscribe(stream pb.RoomGRPC_RoomSubscribeServer) error {
	componentNum, err := getComponentIDFromContext(stream.Context(), c.ComponentIdentifier)
	if err != nil {
		return err
	}
	fmt.Printf("Room Controller : Component%d is registered \n", componentNum)

	myErr := make(chan error)
	msg := make(chan *Supply)
	connectionTable := c.ConnectionTable[componentNum]
	receptionTable := []int32{}

	go func() {
		for {
			recSupply, err := stream.Recv()
			if err == io.EOF {
				// trigger reason unknown
				fmt.Printf("Room Controller ending connections with Component%d due to receinve EOF \n", componentNum)
				myErr <- nil
				return
			}
			if err != nil {
				fmt.Printf("Room Controller ending connections with Component%d due to error, %v \n", componentNum, err)
				myErr <- err
				return
			}
			fmt.Printf("Room Controller : receives %v form Component%d \n", recSupply.ResourceType, componentNum)

			receptionTable = append(receptionTable, recSupply.SenderId)

			for _, neighbourNum := range connectionTable {
				if !slices.Contains(receptionTable, neighbourNum) {
					c.pathways[neighbourNum] <- &Supply{
						senderId:            recSupply.SenderId,
						resourceType:        recSupply.ResourceType,
						resouceAvailability: recSupply.ResouceAvailability,
					}
				}
			}
		}
	}()

	for {
		select {
		case <-myErr:
			return <-myErr
		case rec := <-msg:
			supply := &pb.Supply{
				SenderId:            rec.senderId,
				ResourceType:        rec.resourceType,
				ResouceAvailability: rec.resouceAvailability,
			}
			if err := stream.Send(supply); err != nil {
				fmt.Printf("Room Controller could not reach Component%d, %v \n", componentNum, err)
				myErr <- err
			}
			fmt.Printf("Room Controller : Send %v to Component%d \n", supply, componentNum)
		}
	}
}

func (c *Controller) StartGRPCServer() {
	lis, err := net.Listen("tcp", c.Address)
	if err != nil {
		log.Fatalf("Room Controller : Failed to listen to the controller address at %s, %v", c.Address, err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterRoomGRPCServer(grpcServer, NewController())

	fmt.Printf("Room Controller : Start listening to %s \n", c.Address)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Room Server : Failed to start the server, %v", err)
	}
}
