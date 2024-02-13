package controllerUtil

import (
	"fmt"
	"io"
	"log"
	"net"
	pb "spacegineers_context/protobufs/controller_proto"
	"sync"

	"google.golang.org/grpc"
)

/*
Controller type represents a grpc server that acts as an Event Bus for all components to communicate with each other.
  - pb.UnimplementedRoomGRPCServer: default from grpc implementation
  - serverAddress: is a string that represents the ip address of the Controller server
  - ctxIdentifier: is a string that represents the identifier stored in the context. It is used to retreive componentId through grpc context
  - connectionTables: is a of type map[string][]string that represents the connections between all components. It is obtained through pre-defined yaml file from settings package. The Controller server keep track of these information so that each component can focus on what their own functionalities
  - tunnels: is a sync.Map which is of type asynchronous map[string]chan *Supply. It stores a channel handle to every component's stream so that one can send meesage to the other through channels
*/
type Controller struct {
	pb.UnimplementedControllerGRPCServer
	serverAddress    string
	ctxIdentifier    string
	connectionTables map[string][]string
	tunnels          sync.Map
}

// Controller constructor
func NewController(addr string, identifier string, connections map[string][]string) *Controller {
	return &Controller{
		serverAddress:    addr,
		ctxIdentifier:    identifier,
		connectionTables: connections,
		tunnels:          sync.Map{},
	}
}

// Subscribe method establish server-side streamming between each component client and the Controller server
func (c *Controller) Subscribe(stream pb.ControllerGRPC_SubscribeServer) error {
	componentId, err := getComponentIDFromContext(stream.Context(), c.ctxIdentifier)
	if err != nil {
		return err
	}
	fmt.Printf("Room Controller : Component%s is registered \n", componentId)

	component := NewComponentManager(c.connectionTables[componentId])

	myErr := make(chan error)
	msg := make(chan *pb.Supply)
	c.tunnels.Store(componentId, msg)

	go c.StratReceive(stream, component, componentId, myErr)

	for {
		select {
		case err := <-myErr:
			return err
		case recSupply := <-msg:
			// If the Component hasn't receive a supply from this sender
			component.AddReceptId(recSupply.SenderId, recSupply.ResourceType)
			// Send this sender's message to the Component
			if err := stream.Send(recSupply); err != nil {
				fmt.Printf("Room Controller could not reach Component %s, %v \n", componentId, err)
				myErr <- err
			}
			fmt.Printf("Room Controller : Send {%v} to Component %s \n", recSupply, componentId)

		}
	}
}

func (c *Controller) StratReceive(stream pb.ControllerGRPC_SubscribeServer, component *ComponentManager, componentId string, myErr chan<- error) {
	for {
		recSupply, err := stream.Recv()
		if err == io.EOF {
			// trigger reason unknown
			fmt.Printf("Room Controller ending connections with Component %s due to receinve EOF \n", componentId)
			myErr <- nil
			return
		}
		if err != nil {
			fmt.Printf("Room Controller ending connections with Component %s due to error, %v \n", componentId, err)
			myErr <- err
			return
		}
		fmt.Printf("Room Controller : Component %s receives {%+v} \n", componentId, recSupply)

		recSupply.SenderId = componentId

		for _, neighbourId := range component.neighbourList {
			if component.NotInReceptList(neighbourId, recSupply.ResourceType) {
				channel, ok := c.tunnels.Load(neighbourId)
				if ok {
					channel.(chan *pb.Supply) <- recSupply
				}
			}
		}
	}
}

// StartGRPCServer starts the Controller grpc Server
func (c *Controller) StartGRPCServer() {
	lis, err := net.Listen("tcp", c.serverAddress)
	if err != nil {
		log.Fatalf("Room Controller : Failed to listen to the controller address at %s, %v", c.serverAddress, err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterControllerGRPCServer(grpcServer, c)

	fmt.Printf("Room Controller : Start listening to %s \n", c.serverAddress)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Room Server : Failed to start the server, %v", err)
	}
}
