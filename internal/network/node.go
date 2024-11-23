package network

import (
	"context"
	"fmt"
	"net"
	"sync"

	pb "github.com/21Philip/Auction/internal/grpc"
	"google.golang.org/grpc"
)

type node struct {
	pb.NodeServer
	mu    sync.Mutex
	id    int
	addr  string
	peers map[int]pb.NodeClient // id -> node
	clock *vectorClock
	srv   *grpc.Server // temporary
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

func (n *node) StartNode() {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", n.addr)
	if err != nil {
		fmt.Printf("Unable to start connection to server: %v\n", err)
	}
	fmt.Printf("Node %d listening at %v\n", n.id, listener.Addr())

	pb.RegisterNodeServer(grpcServer, n)
	n.srv = grpcServer

	if grpcServer.Serve(listener) != nil {
		fmt.Printf("Failed to serve: %v\n", err)
	}

	fmt.Printf("Node %d stopped!\n", n.id)
}

func (n *node) Bid(ctx context.Context, in *pb.Amount) (*pb.Ack, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	fmt.Printf("I am node %d and have recived your bid (bidder %d, amount %d)\n", n.id, in.Bidder, in.Amount)

	return &pb.Ack{Success: true}, nil
}

func (n *node) Result(ctx context.Context, in *pb.Empty) (*pb.Outcome, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	fmt.Printf("I am node %d and have recived your request for result\n", n.id)

	return &pb.Outcome{Winner: 0, BidAmount: 0}, nil
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

func (n *node) TestCall(ctx context.Context, in *pb.Empty) (*pb.Test, error) {
	go n.crash()
	response := fmt.Sprintf("Response from node %d", n.id)
	return &pb.Test{Payload: response}, nil
}

func (n *node) crash() {
	n.srv.GracefulStop()
}
