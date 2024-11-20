package node

import (
	"fmt"
	"net"

	pb "/grpc"

	"google.golang.org/grpc"
)

type Node struct {
	pb.NodeServer
}

func main() {

}

func (n *Node) start() {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", n.address)
	if err != nil {
		fmt.Printf("Unable to start connection to server: %v\n", err)
	}
	fmt.Printf("server listening at %v\n", listener.Addr())

	pb.RegisterNodeServer(grpcServer, n)
	if err := n.registerService(); err != nil {
		fmt.Printf("Node %d failed to register\n", n.id)
	}

	go n.nodeLogic(grpcServer)

	if grpcServer.Serve(listener) != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}
