package network

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"

	pb "github.com/21Philip/Auction/internal/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var wg = sync.WaitGroup{}

const BasePort = 50050

type Network struct {
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

func (nw *Network) StartNetwork() {
	for i := range nw.Nodes {
		wg.Add(1)
		go startNode(strconv.Itoa(i), strconv.Itoa(nw.Size))
	}

	wg.Wait()
	fmt.Printf("Network stopped!\n")
}

func startNode(nodeId string, nodeAmount string) {
	cmd := exec.Command("go", "run", "github.com/21Philip/Auction/internal/network/create-node", nodeId, nodeAmount)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		fmt.Printf("cmd.Start failed: %s", err)
	}

	_, err = cmd.Process.Wait()
	if err != nil {
		fmt.Printf("cmd.Process.Wait failed: %s", err)
	}

	wg.Done()
}
