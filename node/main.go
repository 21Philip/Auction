package main

import (
	"fmt"
	"os"
	"strconv"

	pb "github.com/21Philip/Auction/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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