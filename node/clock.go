package main

type CompareResult int

const (
	HappenedBefore CompareResult = iota
	HappenedAfter
	HappenedConcurrently // events are incomparable
	IsSameEvent
)

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

func (vc *VectorClock) getTimestamp(idx int) int {
	return vc.vector[idx]
}

func (vc *VectorClock) merge(other VectorClock) {
	for idx, timestamp := range other.vector {
		vc.vector[idx] = max(vc.vector[idx], timestamp)
	}
}

// Returns whether an event associated with vector clock "vc"
// happened before, after, or concurrently with another event
// having vector clock "other".
// TODO: What if length not same?
func (vc *VectorClock) compareTo(other VectorClock) CompareResult {
	areEqual := true
	otherIsAhead := true
	otherIsBehind := true

	for idx, timestamp := range vc.vector {
		otherTimestamp := other.getTimestamp(idx)
		if timestamp != otherTimestamp {
			areEqual = false
		}
		if timestamp > otherTimestamp {
			otherIsAhead = false
		}
		if timestamp < otherTimestamp {
			otherIsBehind = false
		}
	}

	// vc = other
	if areEqual {
		return IsSameEvent
	}

	// vc -> other
	if otherIsAhead {
		return HappenedBefore
	}

	// other -> vc
	if otherIsBehind {
		return HappenedAfter
	}

	// vc || other
	return HappenedConcurrently
}
