package main

import nwPkg "github.com/21Philip/Auction/internal/network"

func main() {
	nw := nwPkg.NewNetwork(5)
	nw.StartNetwork()
}
