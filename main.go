package main

import (
	"fmt"

	clPkg "github.com/21Philip/Auction/internal/client"
	nwPkg "github.com/21Philip/Auction/internal/network"
)

func main() {
	network, err := nwPkg.NewNetwork(5)
	if err != nil {
		fmt.Printf("ERROR: Could not create network\n")
		return
	}

	go network.StartNetwork()

	client := clPkg.NewClient(0, network)
	client.StartClient()

	network.StopNetwork()
	fmt.Printf("Program stopped!\n")
}
