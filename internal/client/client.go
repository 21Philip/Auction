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
	id      int           // Client id
	nodeId  int           // Id of Current node/replica directing request to
	network nwPkg.Network // All nodes on network
}

func NewClient(id int, network nwPkg.Network) *Client {
	return &Client{
		id:      id,
		nodeId:  0,
		network: network,
	}
}

func (c *Client) StartClient() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		if c.nodeId == -1 {
			break
		}

		input := scanner.Text()
		c.mu.Lock()

		if input == "test" {
			c.testCall()
			continue
		}

		c.mu.Unlock()
	}
}

func (c *Client) testCall() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	req := &pb.Empty{}
	defer cancel()

	reply, err := c.network.Nodes[c.nodeId].TestCall(ctx, req)
	if err != nil {
		c.changeNode(c.testCall)
		return
	}

	fmt.Printf(reply.Payload)
}

func (c *Client) changeNode(retry func()) {
	fmt.Printf("CLIENT (you): Request to current node timed out. Establishing new connection")

	c.nodeId++
	if c.nodeId < c.network.Size {
		retry()
		return
	}

	c.nodeId = -1
}
