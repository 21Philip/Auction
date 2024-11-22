package network

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	pb "github.com/21Philip/Auction/internal/grpc"
	"google.golang.org/grpc"
)

const (
	initialSleepDuration = 2 * time.Second // Allow other nodes to start at beginning of simulation
	stepTime             = 1 * time.Second // The time between each step/frame of simulation
	crashChance          = 10              // The chance of a node to crash at any step. Its calculated as 1/crashChance
)

type node struct {
	pb.NodeServer
	mu    sync.Mutex
	id    int
	addr  string
	peers map[int]pb.NodeClient // id -> node
	clock *vectorClock
	srv   *grpc.Server
}

func NewNode(id int, addr string, peers map[int]pb.NodeClient) *node {
	return &node{
		id:    id,
		addr:  addr,
		peers: peers,
		clock: newVectorClock(),
		srv:   nil,
	}
}

func (n *node) Start() {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", n.addr)
	if err != nil {
		fmt.Printf("Unable to start connection to server: %v\n", err)
	}
	fmt.Printf("Node %d listening at %v\n", n.id, listener.Addr())

	pb.RegisterNodeServer(grpcServer, n)
	n.srv = grpcServer
	//go n.simulateAuction(grpcServer)

	if grpcServer.Serve(listener) != nil {
		fmt.Printf("Failed to serve: %v\n", err)
	}

	fmt.Printf("Node %d stopped!\n", n.id)
}

/*
func (n *node) simulateAuction(srv *grpc.Server) {
	lastStep := time.Now()
	time.Sleep(initialSleepDuration)

	for {
		if time.Since(lastStep) < stepTime {
			continue
		}
		lastStep = time.Now()

		n.mu.Lock()

		fmt.Printf("Hello from node %d\n", n.id)
		if rand.Intn(10) == 0 {
			srv.Stop()
			break
		}

		n.mu.Unlock()
	}

	fmt.Printf("Simulation of node %d was stopped\n", n.id)
}
*/

func (n *node) TestCall(ctx context.Context, in *pb.Empty) (*pb.Test, error) {
	go n.crash()
	response := fmt.Sprintf("Response from node %d", n.id)
	return &pb.Test{Payload: response}, nil
}

func (n *node) Stop(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	if n.srv == nil { // this is impossible
		return nil, fmt.Errorf("node %d was never started", n.id)
	}

	go func() {
		n.srv.GracefulStop()
	}()

	return &pb.Empty{}, nil
}

func (n *node) crash() {
	n.srv.GracefulStop()
}
