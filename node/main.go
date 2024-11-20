package main

import (
	"context"
	"fmt"
	"net"

	//"os"
	//"strconv"
	"time"

	pb "github.com/21Philip/Auction/grpc"
	"google.golang.org/grpc"
)

type vectorClock map[int32]int32

type Node struct {
	pb.NodeServer
	id          int32
	addr        string
	peers       map[int]*Node // id -> node
	vectorClock vectorClock
}

func NewNode(id int32, addr string, peerIDs []int32) *Node {
	clock := make(map[int32]int32)
	for _, peerID := range peerIDs {
		clock[peerID] = 0
	}
	clock[id] = 0

	return &Node{
		id:   id,
		addr: addr,
		//peers:       make([]*Node, 0),
		vectorClock: clock,
	}
}

var peers = make([]*Node, 0)

func (n *Node) incrementClock() {
	n.vectorClock[n.id]++
}

func (n *Node) mergeClock(recievedClock vectorClock) {
	for id, clock := range recievedClock {
		n.vectorClock[id] = max(n.vectorClock[id], clock)
	}
}

// Compare clocks return values:
// -1 if a happens before b
//
//	0 if a and b are concurrent
//	1 if b happens before a
func compareClocks(a, b vectorClock) int {
	aHappenedBeforeB := false
	bHappenedBeforeA := false

	for id, clockA := range a {
		clockB := b[id]

		//a [10, 1]
		//b [1, 10]

		if clockA < clockB {
			bHappenedBeforeA = true
		} else if clockA > clockB {
			aHappenedBeforeB = true
		}
	}

	if aHappenedBeforeB {
		return -1
	} else if bHappenedBeforeA {
		return 1
	} else {
		return 0
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
	/*
		if len(os.Args) != 3 {
			fmt.Printf("ERROR: A node could not be created. REASON: Given invalid number of arguments")
			return
		}

		nodeId, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("ERROR: A node could not be created. REASON: Given invalid id %s", os.Args[1])
			return
		}

		peerAmount, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("ERROR: Node %d could not be created. REASON: Given invalid amount of peers %s", nodeId, os.Args[2])
			return
		}

		for i := range peerAmount {

		}

		port := 50050 + nodeId
		nodeAddress := ":" + strconv.Itoa(port)
	*/

}
