package petrel

import (
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

var (
	ErrEmptyRing    = errors.New("empty ring")
	ErrNodeNotFound = errors.New("node not found in ring")
)

type Node struct {
	ID   uint64
	Addr string
}

type Ring struct {
	ID      uint64
	Owner   Node
	Handler HandlerFunc
	mu      sync.RWMutex
	members []Node
}

// Lookup finds the node currently responsible for a key
func (r *Ring) Lookup(key uint64) (*Node, error) {
	if r == nil {
		return nil, ErrEmptyRing
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.members) == 0 {
		return nil, ErrEmptyRing
	}

	// Simple linear scan for now
	idx := 0
	for i := range r.members {
		if r.members[i].ID > key {
			idx = i
			break
		}
	}

	node := r.members[idx]
	return &node, nil
}

func (r *Ring) AddNode(node Node) error {
	if r == nil {
		return ErrEmptyRing
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.members {
		if r.members[i] == node {
			// Duplicate
			return nil
		}

		if r.members[i].ID > node.ID {
			// Insert just before the existing member
			r.members = append(r.members, Node{})
			copy(r.members[i+1:], r.members[i:])
			r.members[i] = node
			return nil
		}

	}

	r.members = append(r.members, node)
	return nil
}

func (r *Ring) DelNode(node Node) error {
	if r == nil {
		return ErrEmptyRing
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.members {
		if r.members[i] == node {
			copy(r.members[i:], r.members[i+1:])
			r.members = r.members[:len(r.members)-1]
			return nil
		}
	}

	return nil
}

// Pred returns the d'th predecessor to the supplied node.
func (r *Ring) Pred(from Node, d int) (*Node, error) {
	if r == nil {
		return nil, ErrEmptyRing
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.members) == 0 {
		return nil, ErrEmptyRing
	}

	for i := range r.members {
		if r.members[i] == from {
			node := r.members[(i-d)%len(r.members)]
			return &node, nil
		}
	}
	return nil, ErrNodeNotFound
}

// Succ returns the d'th successor to the supplied node.
func (r *Ring) Succ(from Node, d int) (*Node, error) {
	if r == nil {
		return nil, ErrEmptyRing
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.members) == 0 {
		return nil, ErrEmptyRing
	}

	for i := range r.members {
		if r.members[i] == from {
			node := r.members[(i+d)%len(r.members)]
			return &node, nil
		}
	}
	return nil, ErrNodeNotFound
}

func (r *Ring) Send(msg *Message) (*Message, error) {
	if msg.Target == r.Owner {
		// Send locally
		return r.handleMessage(msg)
	}

	conn, err := net.Dial("tcp", msg.Target.Addr)
	if err != nil {
		panic(err)
	}

	codec := TextCodec{}
	codec.Encode(conn, msg)
	conn.Close()

	// Send over network
	return nil, errors.New("not implemented")
}

func (r *Ring) handleMessage(msg *Message) (*Message, error) {
	switch msg.Purpose {
	case MsgUser:
		log.Printf("got user message")
	default:

	}

	return nil, errors.New("not implemented")

}

type HandlerFunc func(io.Reader, io.Writer) error
