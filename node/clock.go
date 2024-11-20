package main

type vectorClock map[int]int

func (n *Node) incrementClock() {
	n.vectorClock[n.id]++
}

func (n *Node) mergeClock(recievedClock vectorClock) {
	for id, clock := range recievedClock {
		n.vectorClock[id] = max(n.vectorClock[id], clock)
	}
}

// Compare clocks return values:
// -1 if a happens before b
//
//	0 if a and b are concurrent
//	1 if b happens before a
func compareClocks(a, b vectorClock) int {
	equal := true
	lessOrEqual := true

	for id, clockA := range a {
		clockB := b[id]
		if clockA != clockB {
			equal = false
			if clockA > clockB {
				lessOrEqual = false
				break
			}
		}
	}

	if equal {
		return 0 //a = b
	}
	if lessOrEqual {
		return -1 //a <= b
	}

	lessOrEqual = true
	for id, clockB := range b {
		clockA := a[id]
		if clockB > clockA {
			lessOrEqual = false
			break
		}
	}

	if lessOrEqual {
		return 1 // b <= a
	}

	return 0 //a || b
}
