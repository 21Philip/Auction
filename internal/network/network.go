package network

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	pb "github.com/21Philip/Auction/internal/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var wg = sync.WaitGroup{}

const (
	BasePort = 50050
	Timeout  = 200 * time.Millisecond
)

type Network struct { // TODO: Make alias
	Size  int
	Nodes map[int]pb.NodeClient // id -> node
}

func NewNetwork(nodeAmount int) (*Network, error) {
	nw := &Network{
		Size:  nodeAmount,
		Nodes: make(map[int]pb.NodeClient),
	}

	for i := range nodeAmount {
		port := BasePort + i
		addr := ":" + strconv.Itoa(port)

		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("failed to create client at iteration %d", i)
		}

		nw.Nodes[i] = pb.NewNodeClient(conn)
	}

	return nw, nil
}

// Blocks until all node processes exit
func (nw *Network) StartNetwork() {
	for i := range nw.Nodes {
		wg.Add(1)
		go startNode(strconv.Itoa(i), strconv.Itoa(nw.Size))
	}

	wg.Wait()
	fmt.Printf("All nodes in network are offline!\n")
}

func startNode(nodeId string, nodeAmount string) {
	cmd := exec.Command("go", "run", "github.com/21Philip/Auction/internal/network/node", nodeId, nodeAmount)
	cmd.Stdout = os.Stdout

	err := cmd.Start()
	if err != nil {
		fmt.Printf("cmd.Start failed: %s", err)
	}

	_, err = cmd.Process.Wait()
	if err != nil {
		fmt.Printf("cmd.Process.Wait failed: %s", err)
	}

	fmt.Printf("Node %s stopped!\n", nodeId)
	wg.Done()
}

// Exits all node processes
func (nw *Network) StopNetwork() {
	for _, node := range nw.Nodes {
		ctx, cancel := context.WithTimeout(context.Background(), Timeout)
		defer cancel()
		node.Stop(ctx, &pb.Empty{})
	}

	wg.Wait()
	time.Sleep(time.Second)
}
