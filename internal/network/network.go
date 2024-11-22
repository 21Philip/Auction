package network

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

var wg = sync.WaitGroup{}

const basePort = 50050

type Network struct {
	Size  int
	Nodes map[int]string // id -> address
}

func NewNetwork(nodeAmount int) *Network {
	nw := &Network{
		Size:  nodeAmount,
		Nodes: make(map[int]string),
	}

	for i := range nodeAmount {
		port := basePort + i
		address := ":" + strconv.Itoa(port)
		nw.Nodes[i] = address
	}

	return nw
}

func (nw *Network) StartNetwork() {
	for i := range nw.Nodes {
		wg.Add(1)
		go startNode(strconv.Itoa(i), strconv.Itoa(nw.Size))
	}

	wg.Wait()
	fmt.Printf("Server stopped!\n")
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
