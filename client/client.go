package client

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	pb "github.com/21Philip/Auction/grpc"
)

const (
	timeout = 200 * time.Millisecond // timeout for all calls to server
)

type Client struct {
	mu         sync.Mutex
	id         int             // Client id
	curNode    int             // Index of node currently directing API requests to
	knownNodes []pb.NodeClient // All known nodes
}

func NewClient(id int, nodes []pb.NodeClient) *Client {
	return &Client{
		id:         id,
		curNode:    0,
		knownNodes: nodes,
	}
}

func (c *Client) Start() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		if c.curNode == -1 {
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

	reply, err := c.knownNodes[c.curNode].TestCall(ctx, req)
	if err != nil {
		c.changeNode(c.testCall)
		return
	}

	fmt.Printf(reply.Payload)
}

func (c *Client) changeNode(retry func()) {
	fmt.Printf("CLIENT (you): Request to current node timed out. Establishing new connection")

	c.curNode++
	if c.curNode < len(c.knownNodes) {
		retry()
		return
	}

	c.curNode = -1
}
