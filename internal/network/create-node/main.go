package main

import (
	"os"
	"strconv"

	nwPkg "github.com/21Philip/Auction/internal/network"
)

func main() {
	id, _ := strconv.Atoi(os.Args[1]) // Do i really care about error handling atp?
	port := nwPkg.BasePort + id
	addr := ":" + strconv.Itoa(port)

	peerAmount, _ := strconv.Atoi(os.Args[2])
	nw := nwPkg.NewNetwork(peerAmount)

	node := nwPkg.NewNode(id, addr, nw.Nodes)
	node.Start()
}
