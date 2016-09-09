package petrel

import (
	"errors"
	"net"
	"sync"
)

type EDRAMessageType uint8

const (
	EDRAMessageEvent EDRAMessageType = iota
	EDRAMessageEventResponse
	EDRAMessageLookup
	EDRAMessageLookupResponse
	EDRAMessageInsert
	EDRAMessageInsertResponse
)

var (
	ErrMissingSelf = errors.New("instance missing self entry in member list")
)

type address struct {
	id   uint64 // the id of an instance on the ring
	addr net.Addr
}

type storedEvent struct {
	source      address
	messageType EDRAMessageType
	ttl         uint32
}

type msgEvent struct {
	messageType  EDRAMessageType
	source       address
	ttl          uint32
	storedEvents []storedEvent
}

type msgEventResponse struct {
	messageType EDRAMessageType
	source      address
	ttl         uint32
}

type msgLookup struct {
	messageType EDRAMessageType
	source      address
}

type msgLookupResponse struct {
	messageType EDRAMessageType
	source      address
	succ        address
}

type msgInsert struct {
	messageType EDRAMessageType
	source      address
}

type msgInsertResponse struct {
	messageType EDRAMessageType
	source      address
	newMember   address
	memberList  []address
}

type instance struct {
	addr address

	mu           sync.RWMutex
	memberList   []address
	selfIndex    int
	storedEvents []storedEvent
}

// note: locateSelf does no locking and assumes its caller has a read
// or write lock on the member list
func (in *instance) locateSelf() error {
	for i := range in.memberList {
		if in.memberList[i] == in.addr {
			in.selfIndex = i
			return nil
		}
	}

	return ErrMissingSelf
}

// delMember adds the given address to the member list.
func (in *instance) addMember(addr address) error {
	in.mu.Lock()
	defer in.mu.Unlock()

	for i := range in.memberList {
		if in.memberList[i] == addr {
			// Duplicate
			return nil
		}

		if in.memberList[i].id > addr.id {
			// Insert just before the existing member
			in.memberList = append(in.memberList, address{})
			copy(in.memberList[i+1:], in.memberList[i:])
			in.memberList[i] = addr
			return in.locateSelf()
		}

	}

	in.memberList = append(in.memberList, addr)
	return in.locateSelf()
}

// delMember removes the given address from the member list.
func (in *instance) delMember(addr address) error {
	in.mu.Lock()
	defer in.mu.Unlock()

	for i := range in.memberList {
		if in.memberList[i] == addr {
			copy(in.memberList[i:], in.memberList[i+1:])
			in.memberList = in.memberList[:len(in.memberList)-1]
			return in.locateSelf()
		}
	}
	return nil
}

// pred returns the d'th predecessor to the current instance.
func (in *instance) pred(d int) address {
	in.mu.RLock()
	defer in.mu.RUnlock()
	return in.memberList[(in.selfIndex-d)%len(in.memberList)]
}

// succ returns the d'th successor to the current instance.
func (in *instance) succ(d int) address {
	in.mu.RLock()
	defer in.mu.RUnlock()
	return in.memberList[(in.selfIndex+d)%len(in.memberList)]
}

func (in *instance) storeEvent(e storedEvent) {
	in.mu.Lock()
	defer in.mu.Unlock()
	in.storedEvents = append(in.storedEvents, e)
}

func (in *instance) resetEvents() {
	in.mu.Lock()
	defer in.mu.Unlock()
	in.storedEvents = in.storedEvents[:0]
}

func log2ceil(n int) int {
	if n == 0 {
		return 0
	}

	x := 0

	// Check if it's a power of 2
	if (n & (n - 1)) != 0 {
		x++
	}

	// Count bits
	for n > 1 {
		n >>= 1
		x++
	}

	return x
}

func pow2(n int) int {
	x := 1
	for n > 0 {
		x <<= 1
		n--
	}
	return x
}

func (in *instance) thetaFn() {
	in.mu.RLock()

	if len(in.storedEvents) == 0 || len(in.memberList) <= 1 {
		in.mu.RUnlock()
		return
	}

	distance := log2ceil(len(in.memberList))
	in.mu.RUnlock()

	for ttl := 0; ttl < distance; ttl++ {
		target := in.succ(pow2(ttl))

		in.transmitEvents(ttl, in.storedEvents, target)
	}
	in.resetEvents()
}

func (in *instance) transmitEvents(ttl int, events []storedEvent, target address) {

}
