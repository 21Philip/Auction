# Auction
Distributed auction system

## To get started
1. Clone the repository
2. Configure the amount of nodes/replicas in Auction/main.go
   ```go
    network, err := nwPkg.NewNetwork(<number of nodes>)
    ```
   The default is 5. You can change this value to anything greater than 3.
3. In the root of the repo, start the application with:
    ```bash
    go run .
    ```
    
## Use
Once the program is started you can type any of the following commands to interact with the system:
- 'bid \<amount\>': Places a bid in the system
- 'result': Displays the current highest bid and the id of the client that placed that bid.  
  NOTE: Since you are the only client connected to the system, the highest bid will always be you (clientId = 0).
  The one exception is if no bids have been placed, in which case the highest bidder is clientId = -1.
- 'quit': Exits all remaining node processes and the client.
- 'kill \<nodeId\>': Immediately crashes the node with the given id. Used to demonstrate the fault tolerance.  
  Note: The nodes are given id's 0, 1, 2, ..., nNodes - 1. The client will connect to node 0 initially, and increment
  by one every time a request times out. If the client runs out of nodes, it will inform you to use 'quit' to exit the program.
  If half or more of the nodes are killed, the remaining nodes will stop responding and the auction cease to function.

## Requirements
The program is build using go version 1.23.0

## Implementation notes
We assume a network that has reliable, ordered message transport, where transmissions to non-failed nodes complete within a known time-limit.
