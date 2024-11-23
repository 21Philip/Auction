package network

type compareResult int

const (
	HappenedBefore compareResult = iota
	HappenedAfter
	HappenedConcurrently // events are incomparable
	IsSameEvent
)

type vectorClock struct {
	vector map[int]int // idx -> timestamp
}

func newVectorClock() *vectorClock {
	return &vectorClock{
		vector: make(map[int]int),
	}
}

func (vc *vectorClock) incrementTimestamp(idx int) {
	vc.vector[idx]++
}

func (vc *vectorClock) getTimestamp(idx int) int {
	return vc.vector[idx]
}

func (vc *vectorClock) merge(other vectorClock) {
	for idx, timestamp := range other.vector {
		vc.vector[idx] = max(vc.vector[idx], timestamp)
	}
}

// Returns whether an event associated with vector clock "vc"
// happened before, after, or concurrently with another event
// having vector clock "other".
// TODO: What if length not same?
func (vc *vectorClock) compareTo(other vectorClock) compareResult {
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
