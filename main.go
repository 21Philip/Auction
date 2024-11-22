package main

import (
	"fmt"

	nwPkg "github.com/21Philip/Auction/internal/network"
)

func main() {
	nw, err := nwPkg.NewNetwork(5)
	if err != nil {
		fmt.Printf("ERROR: Could not create network")
		return
	}
	nw.StartNetwork()
}
