package petrel

import (
	"errors"
	"net"
)

type Node struct {
	ID   uint64 // the id of an instance on the ring
	Addr net.Addr
}

type NodeSet struct {
	Primary  *Node
	Replicas []*Node
}

type Ring struct {
	mu      sync.RWMutex
	members []Node
}

// Lookup finds the nodeset currently responsible for a key
func (r *Ring) Lookup(key uint64) (NodeSet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for i := range r.members {
		if in.memberList[i] == in.addr {
			in.selfIndex = i
			return nil
		}
	}

	return NodeSet{}, errors.New("Not implemented")
}
