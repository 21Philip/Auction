package server

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

var wg = sync.WaitGroup{}

func StartServer(nodeAmount int) {
	for i := range nodeAmount {
		wg.Add(1)
		go startNode(strconv.Itoa(i), strconv.Itoa(nodeAmount))
	}

	wg.Wait()
	fmt.Printf("Server stopped!\n")
}

func startNode(nodeId string, nodeAmount string) {
	cmd := exec.Command("go", "run", "github.com/21Philip/Auction/internal/server/create-node", nodeId, nodeAmount)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		fmt.Printf("cmd.Start failed: %s", err)
	}

	_, err = cmd.Process.Wait()
	if err != nil {
		fmt.Printf("cmd.Process.Wait failed: %s", err)
	}

	wg.Done()
}
