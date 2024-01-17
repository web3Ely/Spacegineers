package main

import (
	"context"
	"fmt"
	"log"
	"net"
	pb "spacegineers_context/protobufs/context_proto"

	"google.golang.org/grpc"
)

// ----------------------------------------------
// Worker pool Declaration and implementation
// ----------------------------------------------
// An InJob is just a damage comming from other player
// An OutJob is context api send to designated room
type InJob struct {
	Value int32
}
type OutJob struct {
	Value int32
}

// JobQueueManager represents a manager struct that manges all type of job channels
// Every time another player initiates a grpc call to this context API, it queue them in JobInChannel
// Once context api finish calculating the impack of the damage done to rooms, it pack the result to JobOutChannel to be sent to rooms
// For now, only one job can be queued in each group
type JobQueueManager struct {
	JobInChannel  chan InJob
	JobOutChannel chan OutJob
}

// Returns a pointer to a JobQueueManager
// Called once
func NewJobQueueManager() *JobQueueManager {
	jobInChannel := make(chan InJob)
	jobOutChannel := make(chan OutJob)
	return &JobQueueManager{JobInChannel: jobInChannel, JobOutChannel: jobOutChannel}
}

// Start method binds to a JobQueueManager pointer, it spins off workers goroutine to handler jobs
func (j *JobQueueManager) Start(maxWorker int) {
	for i := 0; i < maxWorker; i++ {
		go Worker(j.JobInChannel, j.JobOutChannel)
	}
}

// A worker listens to the JobInChannel for incoming Jobs
// Once it receives a job, it start its handling process
func Worker(jobInChannel <-chan InJob, jobOutChannel chan<- OutJob) {
	for inJob := range jobInChannel {
		log.Printf("received a InJob %d", inJob)
		// time.Sleep(time.Second * 3)
		outJob := OutJob{Value: 10}
		jobOutChannel <- outJob
		log.Println("Finshed sending the InJob to a Room")
	}
}

// ----------------------------------------------
// End
// ----------------------------------------------

// ----------------------------------------------
// GRPC Server implementation
// ----------------------------------------------
// Define GRPC Server handler
type Server struct {
	pb.UnimplementedSpaceContextServer
	JobQueueManager *JobQueueManager
}

// Hit method is activated when another player fires to this ship
// It is blocking, a respond is blocked when all workers are working, that means the receiver will not receive a respond on time
// (Change) Might want to store incoming request so that the http communication is not blocked
func (s *Server) Hit(ctx context.Context, req *pb.Damage) (*pb.Empty, error) {
	log.Printf("Received a damage: %v", req.Damage)
	job := InJob{Value: req.Damage}
	s.JobQueueManager.JobInChannel <- job
	log.Printf("Pass")
	return &pb.Empty{}, nil
}

func (s *Server) RoomRegister(req *pb.Empty, stream pb.SpaceContext_RoomRegisterServer) error {
	log.Printf("Received a room registration request")
	// for i := 0; i < 3; i++ {
	for outJob := range s.JobQueueManager.JobOutChannel {
		log.Println("Send impact result to a room")
		if err := stream.Send(&pb.Damage{Damage: outJob.Value}); err != nil {
			log.Printf("Room register error: %v \n", err)
			return err
		}
	}
	return nil
}

// ----------------------------------------------
// End
// ----------------------------------------------

func main() {
	//const change to enviornment variable
	const (
		port      = 8080
		maxWorker = 3
	)

	log.Println("Initializing the Context")

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("Failed to listen to the port %d: %v", port, err)
	}
	log.Printf("Listening to the localhost:%d", port)

	grpcServer := grpc.NewServer()
	log.Println("Starting the gRPC Server")

	jobQueue := NewJobQueueManager()
	jobQueue.Start(maxWorker)

	server := &Server{JobQueueManager: jobQueue}
	pb.RegisterSpaceContextServer(grpcServer, server)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}

}

// improvements:
// 1. Make JobInChannel/JobOutChannel buffered channel so that it does not lock the traffic
// 2. Set context deadline for a player when they call grpc functions to avoid dead lock
