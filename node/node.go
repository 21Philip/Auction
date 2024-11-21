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
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Node struct {
	pb.NodeServer
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
	go n.nodeLogic()

	if grpcServer.Serve(listener) != nil {
		fmt.Printf("Failed to serve: %v\n", err)
	}
}

func (n *Node) nodeLogic() {
	time.Sleep(2 * time.Second)
	for id, peer := range n.peers {
		if id == n.id+1 {
			peer.TestCall(context.Background(), &pb.Empty{})
		}
	}
}

func (n *Node) TestCall(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	fmt.Printf("I am node %d\n", n.id)

	md, _ := metadata.FromIncomingContext(ctx)
	for k, v := range md {
		fmt.Printf("%s:\n", k)
		for _, s := range v {
			fmt.Printf("   %s\n", s)
		}
	}

	fmt.Printf("\n")
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

		conn, err := grpc.NewClient(peerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("ERROR: Node %d could not connect to %s: %v\n", node.id, addr, err)
			continue
		}

		node.peers[i] = pb.NewNodeClient(conn)
	}

	node.start()
}
