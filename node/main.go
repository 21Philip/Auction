package main

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	pb "github.com/21Philip/Auction/grpc"
	"google.golang.org/grpc"
)

type Node struct {
	pb.NodeServer
	id   int32
	addr string
}

func NewNode(id int32, addr string) *Node {
	return &Node{
		id:   id,
		addr: addr,
	}
}

var peers = make([]*Node, 0)

func (n *Node) start() {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", n.addr)
	if err != nil {
		fmt.Printf("Unable to start connection to server: %v\n", err)
	}
	fmt.Printf("server listening at %v\n", listener.Addr())

	pb.RegisterNodeServer(grpcServer, n)

	peers = append(peers, n)
	go n.nodeLogic()

	if grpcServer.Serve(listener) != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}

func (n *Node) nodeLogic() {
	time.Sleep(2 * time.Second)
	for _, peer := range peers {
		if peer.id != n.id {
			peer.TestCall(context.Background(), &pb.Empty{})
		}
	}
}

func (n *Node) TestCall(_ context.Context, in *pb.Empty) (*pb.Empty, error) {
	fmt.Printf("I am node %d\n", n.id)
	return &pb.Empty{}, nil
}

func main() {
	n0 := NewNode(0, ":50050")
	n1 := NewNode(1, ":50051")

	wg := sync.WaitGroup{}
	wg.Add(1)

	go n0.start()
	go n1.start()

	wg.Wait()
}
