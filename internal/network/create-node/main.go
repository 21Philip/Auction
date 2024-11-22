package main

import (
	"fmt"
	"os"
	"strconv"

	nwPkg "github.com/21Philip/Auction/internal/network"
)

func main() {
	id, _ := strconv.Atoi(os.Args[1]) // Do i really care about error handling here?
	port := nwPkg.BasePort + id
	addr := ":" + strconv.Itoa(port)

	peerAmount, _ := strconv.Atoi(os.Args[2]) // no
	nw, err := nwPkg.NewNetwork(peerAmount)
	if err != nil {
		fmt.Printf("ERROR: Node %d could not be created due to network error: %v", id, err)
		return
	}

	node := nwPkg.NewNode(id, addr, nw.Nodes)
	node.Start()
}
