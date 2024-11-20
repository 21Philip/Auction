package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	pb "github.com/21Philip/Auction/grpc"
	"google.golang.org/grpc"
)

type Node struct {
	pb.NodeServer
	id          int
	addr        string
	peers       map[int]pb.NodeClient // id -> node
	vectorClock vectorClock
}

func NewNode(id int, addr string) *Node {
	return &Node{
		id:          id,
		addr:        addr,
		peers:       make(map[int]pb.NodeClient),
		vectorClock: make(map[int]int),
	}
}

func (n *Node) start() {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", n.addr)
	if err != nil {
		fmt.Printf("Unable to start connection to server: %v\n", err)
	}
	fmt.Printf("server listening at %v\n", listener.Addr())

	pb.RegisterNodeServer(grpcServer, n)
	go n.nodeLogic()

	if grpcServer.Serve(listener) != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}

func (n *Node) nodeLogic() {
	time.Sleep(2 * time.Second)
	for id, peer := range n.peers {
		if id != n.id {
			peer.TestCall(context.Background(), &pb.Empty{})
		}
	}
}

func (n *Node) TestCall(_ context.Context, in *pb.Empty) (*pb.Empty, error) {
	fmt.Printf("I am node %d\n", n.id)
	return &pb.Empty{}, nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("ERROR: A node could not be created. REASON: Invalid number of arguments")
		return
	}

	nodeId, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("ERROR: A node could not be created. REASON: Invalid id %s", os.Args[1])
		return
	}

	peerAmount, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("ERROR: Node %d could not be created. REASON: Invalid amount of peers %s", nodeId, os.Args[2])
		return
	}

	basePort := 50050
	port := basePort + nodeId
	addr := ":" + strconv.Itoa(port)

	node := NewNode(nodeId, addr)

	for i := range peerAmount {
		peerAddr := ":" + strconv.Itoa(basePort+i)

		conn, err := grpc.NewClient(peerAddr)
		if err != nil {
			fmt.Printf("ERROR: Node %d could not connect to %s: %v\n", node.id, addr, err)
			continue
		}

		node.peers[i] = pb.NewNodeClient(conn)
		node.vectorClock[i] = 0
	}

	node.start()
}
