package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	pb "github.com/21Philip/Auction/grpc"
	"google.golang.org/grpc"
)

const (
	initialSleepDuration = 2 * time.Second // Allow other nodes to start at beginning of simulation
	stepTime             = 1 * time.Second // The time between each step/frame of simulation
	crashChance          = 10              // The chance of a node to crash at any step. Its calculated as 1/crashChance
)

type Node struct {
	pb.NodeServer
	mu    sync.Mutex
	id    int
	addr  string
	peers map[int]pb.NodeClient // id -> node
	clock *VectorClock
}

func NewNode(id int, addr string) *Node {
	return &Node{
		id:    id,
		addr:  addr,
		peers: make(map[int]pb.NodeClient),
		clock: NewVectorClock(),
	}
}

func (n *Node) start() {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", n.addr)
	if err != nil {
		fmt.Printf("Unable to start connection to server: %v\n", err)
	}
	fmt.Printf("Node %d listening at %v\n", n.id, listener.Addr())

	pb.RegisterNodeServer(grpcServer, n)
	go n.simulateAuction(grpcServer)

	if grpcServer.Serve(listener) != nil {
		fmt.Printf("Failed to serve: %v\n", err)
	}

	fmt.Printf("Node %d was killed\n", n.id)
}

func (n *Node) simulateAuction(srv *grpc.Server) {
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

func (n *Node) TestCall(ctx context.Context, in *pb.Empty) (*pb.Test, error) {
	fmt.Printf("I am node %d\n", n.id)
	return &pb.Test{Payload: "response"}, nil
}
