package main

import (
	"fmt"
	"os"
	"strconv"

	pb "github.com/21Philip/Auction/internal/grpc"
	nwPkg "github.com/21Philip/Auction/internal/network"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	nodeId, _ := strconv.Atoi(os.Args[1]) // Do i really care about error handling atp?
	peerAmount, _ := strconv.Atoi(os.Args[2])

	nw := nwPkg.NewNetwork(peerAmount)
	peers := make(map[int]pb.NodeClient)

	for id, addr := range nw.Nodes {

		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("ERROR: Node %d could not connect to %s: %v\n", nodeId, addr, err)
			continue
		}

		peers[id] = pb.NewNodeClient(conn)
	}

	nodeAddr := nw.Nodes[nodeId]
	node := nwPkg.NewNode(nodeId, nodeAddr, peers)
	node.Start()
}
