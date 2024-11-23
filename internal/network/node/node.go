package main

import (
	"context"
	"fmt"
	"net"
	"sync"

	pb "github.com/21Philip/Auction/internal/grpc"
	nwPkg "github.com/21Philip/Auction/internal/network"
	"google.golang.org/grpc"
)

type node struct {
	pb.NodeServer
	mu         sync.Mutex
	id         int
	addr       string
	network    *nwPkg.Network
	highestBid *pb.Amount
	srv        *grpc.Server // temporary TODO: Remove
}

func newNode(id int, addr string, network *nwPkg.Network) *node {
	return &node{
		id:         id,
		addr:       addr,
		network:    network,
		highestBid: &pb.Amount{Bidder: -1, Amount: 0},
		srv:        nil,
	}
}

func (n *node) startNode() {
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

func (n *node) getMajorityApproval(bid *pb.Amount) bool {
	majority := n.network.Size / 2 // Already has approvement from self
	approvals := 0

	ctx, cancel := context.WithTimeout(context.Background(), nwPkg.Timeout)
	defer cancel()

	for peerId, peer := range n.network.Nodes {
		go func() {
			if peerId != n.id { // TODO: Consider looping over self aswell
				return
			}

			reply, err := peer.VerifyBid(ctx, bid)
			if err == nil && reply.Success {
				approvals++
			}
		}()
	}

	return approvals >= majority
}

func (n *node) VerifyBid(ctx context.Context, in *pb.Amount) (*pb.Ack, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.highestBid.Amount > in.Amount {
		return &pb.Ack{Success: false}, nil
	}

	n.highestBid = in
	return &pb.Ack{Success: true}, nil
}

func (n *node) Bid(ctx context.Context, in *pb.Amount) (*pb.Ack, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	fmt.Printf("Node %d: Recived bid (bidder %d, amount %d)\n", n.id, in.Bidder, in.Amount)

	if n.highestBid.Amount > in.Amount {
		return &pb.Ack{Success: false}, nil
	}

	// check if majority approve incoming bid
	if !n.getMajorityApproval(in) {
		return &pb.Ack{Success: false}, nil
	}

	// If so, change local leadingBid
	n.highestBid = in
	return &pb.Ack{Success: true}, nil
}

func (n *node) Result(ctx context.Context, in *pb.Empty) (*pb.Outcome, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	fmt.Printf("Node %d: Recived request for result\n", n.id)

	return &pb.Outcome{HighestBid: n.highestBid}, nil
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
