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
}

func newNode(id int, addr string, network *nwPkg.Network) *node {
	return &node{
		id:         id,
		addr:       addr,
		network:    network,
		highestBid: &pb.Amount{Bidder: -1, Amount: 0},
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

	if grpcServer.Serve(listener) != nil {
		fmt.Printf("Failed to serve: %v\n", err)
	}
}

// Returns true if majority of network approves given bid.
// Assumes bid is already approved by caller
func (n *node) getMajorityApproval(bid *pb.Amount) bool {
	majority := n.network.Size / 2
	approvals := make(chan bool, n.network.Size)

	ctx, cancel := context.WithTimeout(context.Background(), nwPkg.Timeout)
	defer cancel()

	calls := sync.WaitGroup{}

	for peerId, peer := range n.network.Nodes {
		if peerId == n.id { // TODO: Consider looping over self aswell. locks acting wierd
			continue
		}
		calls.Add(1)

		go func() {
			reply, err := peer.VerifyBid(ctx, bid)
			if err == nil && reply.Success {
				approvals <- true
			}
			calls.Done()
		}()
	}

	calls.Wait()
	return len(approvals) >= majority
}

func (n *node) VerifyBid(ctx context.Context, in *pb.Amount) (*pb.Ack, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.highestBid.Amount > in.Amount {
		return &pb.Ack{Success: false}, nil
	}
	if n.highestBid.Amount == in.Amount && n.highestBid.Bidder != in.Bidder {
		return &pb.Ack{Success: false}, nil
	}

	n.highestBid = in
	return &pb.Ack{Success: true}, nil
}

func (n *node) Bid(ctx context.Context, in *pb.Amount) (*pb.Ack, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.highestBid.Amount > in.Amount {
		return &pb.Ack{Success: false}, nil
	}

	if n.highestBid.Amount == in.Amount && n.highestBid.Bidder != in.Bidder {
		return &pb.Ack{Success: false}, nil
	}

	if !n.getMajorityApproval(in) {
		return nil, fmt.Errorf("splitbrain")
	}

	n.highestBid = in
	return &pb.Ack{Success: true}, nil
}

func (n *node) Result(ctx context.Context, in *pb.Empty) (*pb.Outcome, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if !n.getMajorityApproval(n.highestBid) {
		return nil, fmt.Errorf("splitbrain")
	}

	return &pb.Outcome{HighestBid: n.highestBid}, nil
}

// Synthesizes a crash
func (n *node) Stop(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	null := n.network.Nodes[-1]
	null.Stop(ctx, in)
	return &pb.Empty{}, nil
}
