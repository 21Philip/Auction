package client

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	pb "github.com/21Philip/Auction/internal/grpc"
	nwPkg "github.com/21Philip/Auction/internal/network"
)

const (
	timeout = 200 * time.Millisecond // timeout for all calls to server
)

type Client struct {
	mu      sync.Mutex
	id      int            // Client id
	nodeId  int            // Id of Current node/replica directing request to
	network *nwPkg.Network // All nodes on network
}

func NewClient(id int, network *nwPkg.Network) *Client {
	return &Client{
		id:      id,
		nodeId:  0,
		network: network,
	}
}

// Blocks until user types 'quit' or recieves interrupt signal
func (c *Client) StartClient() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input := scanner.Text()
		c.mu.Lock()

		if input == "test" {
			c.testCall()
		}

		if input == "quit" {
			break
		}

		c.mu.Unlock()
	}

	fmt.Println("Client stopped!")
}

func (c *Client) testCall() {
	curNode := c.network.Nodes[c.nodeId]
	if curNode == nil {
		fmt.Printf("CLIENT (you): Seems all nodes are unavailable. Consider using 'quit' to exit program\n")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	req := &pb.Empty{}
	defer cancel()

	reply, err := curNode.TestCall(ctx, req)
	if err != nil {
		c.changeNode(c.testCall)
		return
	}

	fmt.Printf("%s\n", reply.Payload)
}

// Very simple for now. Should consider better logic at some point
func (c *Client) changeNode(retry func()) {
	fmt.Printf("CLIENT (you): Request to current node timed out. Establishing new connection\n")
	c.nodeId++
	retry()
}
