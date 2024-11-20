package main

type VectorClock struct {
	vector map[int]int // idx -> timestamp
}

func NewVectorClock() *VectorClock {
	return &VectorClock{
		vector: make(map[int]int),
	}
}

func (vc *VectorClock) incrementTimestamp(idx int) {
	vc.vector[idx]++
}

func (vc *VectorClock) merge(other VectorClock) {
	for idx, timestamp := range other.vector {
		vc.vector[idx] = max(vc.vector[idx], timestamp)
	}
}

func (vc *VectorClock) compareTo(other VectorClock) int {
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

// Compare clocks return values:
// -1 if a happens before b
//	0 if a and b are concurrent
//	1 if b happens before a
