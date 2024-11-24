package client

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	pb "github.com/21Philip/Auction/internal/grpc"
	nwPkg "github.com/21Philip/Auction/internal/network"
	"google.golang.org/grpc"
)

const (
	timeout = 2 * nwPkg.Timeout // timeout for all calls to server
)

type Client struct {
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

		if input[0] == "bid" {
			c.bid(input)
		}

		if input[0] == "result" {
			c.result()
		}

		if input[0] == "quit" {
			break
		}

		if input[0] == "kill" {
			c.killNode(input)
		}
	}

	fmt.Println("Client stopped!")
}

func makeCall[In any, Out any](c *Client, call func(pb.NodeClient, context.Context, In, ...grpc.CallOption) (Out, error), req In) (Out, error) {
	var reply Out

	curNode := c.network.Nodes[c.nodeId]
	if curNode == nil { // TODO: Consider cycling nodes
		fmt.Printf("Client (you): Seems all nodes are unavailable. Consider using 'quit' to exit program\n")
		return reply, fmt.Errorf("no more nodes")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	reply, err := call(curNode, ctx, req)
	if err != nil {
		fmt.Printf("Client (you): Request to node %d timed out. Establishing new connection\n", c.nodeId)
		c.nodeId++
		return makeCall(c, call, req)
	}

	return reply, nil
}

func (c *Client) bid(input []string) {
	if len(input) != 2 {
		fmt.Printf("Client (you): Incorrect arguments to place bid. Correct use 'bid <amount>'\n")
		return
	}
	bidAmount, err := strconv.Atoi(input[1])
	if err != nil {
		fmt.Printf("Client (you): Cannot convert '%s' to int. Correct use 'bid <amount>'\n", input[1])
		return
	}

	req := &pb.Amount{
		Bidder: int32(c.id),
		Amount: int32(bidAmount),
	}

	reply, err := makeCall(c, pb.NodeClient.Bid, req)
	if err != nil {
		fmt.Printf("Client (you): Something went wrong. Bid not placed.\n")
		return
	}

	if reply.Success {
		fmt.Printf("Client (you): Success! Your bid of %d$ has been placed.\n", bidAmount)
		return
	}

	fmt.Println("Client (you): Bid too low")
	c.result()
}

func (c *Client) result() {
	req := &pb.Empty{}
	reply, err := makeCall(c, pb.NodeClient.Result, req)
	if err != nil {
		return
	}

	fmt.Printf("Client (you): Highest bidder: %d, bid %d$\n", reply.HighestBid.Bidder, reply.HighestBid.Amount)
}

// For testing
func (c *Client) killNode(input []string) {
	if len(input) != 2 {
		return
	}
	nodeId, err := strconv.Atoi(input[1])
	if err != nil {
		return
	}

	req := &pb.Empty{}
	NodeToKill := c.network.Nodes[nodeId]

	if NodeToKill != nil {
		NodeToKill.Stop(context.Background(), req)
	}
}
