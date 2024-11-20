package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

var wg = sync.WaitGroup{}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("FATAL: Incorrect number of arguments\n")
		return
	}

	nodeAmountStr := os.Args[1]
	nodeAmountInt, err := strconv.Atoi(nodeAmountStr)
	if err != nil {
		fmt.Printf("FATAL: Could not convert argument %s to int\n", nodeAmountStr)
		return
	}

	for i := range nodeAmountInt {
		wg.Add(1)
		go runNode(strconv.Itoa(i), nodeAmountStr)
	}

	wg.Wait()
	fmt.Printf("Finished!\n")
}

func runNode(nodeId string, nodeAmount string) {
	cmd := exec.Command("go", "run", "./node", nodeId, nodeAmount)
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
