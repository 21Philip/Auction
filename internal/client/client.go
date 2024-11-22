package client

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
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
		input := strings.Split(scanner.Text(), " ")
		c.mu.Lock()

		if input[0] == "status" {
			c.status()
		}

		if input[0] == "bid" {
			c.bid(input)
		}

		if input[0] == "quit" {
			break
		}

		if input[0] == "test" {
			c.testCall()
		}

		c.mu.Unlock()
	}

	fmt.Println("Client stopped!")
}

func (c *Client) status() {
	curNode := c.network.Nodes[c.nodeId]
	if curNode == nil {
		fmt.Printf("CLIENT (you): Seems all nodes are unavailable. Consider using 'quit' to exit program\n")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	req := &pb.Empty{}
	defer cancel()

	reply, err := curNode.Result(ctx, req)
	if err != nil {
		c.changeNode(func() { c.status() })
		return
	}

	fmt.Printf("Current winner: Client %d, bid %d\n", reply.Winner, reply.BidAmount)
}

func (c *Client) bid(input []string) {
	curNode := c.network.Nodes[c.nodeId]
	if curNode == nil {
		fmt.Printf("CLIENT (you): Seems all nodes are unavailable. Consider using 'quit' to exit program\n")
		return
	}

	if len(input) != 2 {
		fmt.Printf("CLIENT (you): Incorrect arguments to place bid. Correct use 'bid <amount>'")
		return
	}

	bidAmount, err := strconv.Atoi(input[1])
	if err != nil {
		fmt.Printf("CLIENT (you): Cannot convert %s to int. Correct use 'bid <amount>'", input[1])
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	req := &pb.Amount{
		Bidder: int32(c.id),
		Amount: int32(bidAmount),
	}
	defer cancel()

	reply, err := curNode.Bid(ctx, req)
	if err != nil {
		c.changeNode(func() { c.bid(input) })
		return
	}

	fmt.Printf("%v\n", reply.Success)
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
		c.changeNode(func() { c.testCall() })
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
