package petrel

import (
	"reflect"
	"testing"
)

func TestLookup(t *testing.T) {
	testCases := []struct {
		ring   *Ring
		key    uint64
		nodeID uint64
		err    error
	}{

		{
			ring: nil,
			err:  ErrEmptyRing,
		},

		{
			ring: &Ring{
				members: []Node{},
			},
			err: ErrEmptyRing,
		},

		{
			ring: &Ring{
				members: []Node{
					Node{ID: 100},
				},
			},
			key:    99,
			nodeID: 100,
		},

		{
			ring: &Ring{
				members: []Node{
					Node{ID: 100},
				},
			},
			key:    100,
			nodeID: 100,
		},

		{
			ring: &Ring{
				members: []Node{
					Node{ID: 100},
				},
			},
			key:    101,
			nodeID: 100,
		},

		{
			ring: &Ring{
				members: []Node{
					Node{ID: 30},
					Node{ID: 100},
					Node{ID: 300},
				},
			},
			key:    10,
			nodeID: 30,
		},

		{
			ring: &Ring{
				members: []Node{
					Node{ID: 30},
					Node{ID: 100},
					Node{ID: 300},
				},
			},
			key:    99,
			nodeID: 100,
		},

		{
			ring: &Ring{
				members: []Node{
					Node{ID: 30},
					Node{ID: 100},
					Node{ID: 300},
				},
			},
			key:    101,
			nodeID: 300,
		},
	}

	for i, tc := range testCases {
		node, err := tc.ring.Lookup(tc.key)
		if err != tc.err {
			t.Errorf("%d: got error %v, wanted %v", i, err, tc.err)
		}

		if tc.err != nil {
			continue
		}

		if node.ID != tc.nodeID {
			t.Errorf("%d: got node %v, wanted %v", i, node.ID, tc.nodeID)
		}
	}

}

func TestAddNode(t *testing.T) {

	node100 := Node{ID: 100, Addr: "1.1.1.1:1000"}
	node200 := Node{ID: 200, Addr: "1.1.1.2:1000"}
	node300 := Node{ID: 300, Addr: "1.1.1.3:1000"}

	testCases := []struct {
		ring     *Ring
		node     Node
		expected *Ring
		err      error
	}{
		{
			ring:     &Ring{members: []Node{}},
			node:     node100,
			expected: &Ring{members: []Node{node100}},
		},

		{
			ring:     &Ring{members: []Node{node100}},
			node:     node200,
			expected: &Ring{members: []Node{node100, node200}},
		},

		{
			ring:     &Ring{members: []Node{node100}},
			node:     node100,
			expected: &Ring{members: []Node{node100}},
		},

		{
			ring:     &Ring{members: []Node{node200}},
			node:     node100,
			expected: &Ring{members: []Node{node100, node200}},
		},

		{
			ring:     &Ring{members: []Node{node100, node300}},
			node:     node200,
			expected: &Ring{members: []Node{node100, node200, node300}},
		},
	}

	for i, tc := range testCases {
		err := tc.ring.AddNode(tc.node)

		if err != tc.err {
			t.Errorf("%d: got error %v, wanted %v", i, err, tc.err)
		}
		if tc.err != nil {
			continue
		}

		if !reflect.DeepEqual(tc.ring.members, tc.expected.members) {
			t.Errorf("%d: got %+v, wanted %+v", i, tc.ring.members, tc.expected.members)
		}
	}
}

func TestDelNode(t *testing.T) {

	node100 := Node{ID: 100, Addr: "1.1.1.1:1000"}
	node200 := Node{ID: 200, Addr: "1.1.1.2:1000"}
	node300 := Node{ID: 300, Addr: "1.1.1.3:1000"}

	testCases := []struct {
		ring     *Ring
		node     Node
		expected *Ring
		err      error
	}{
		{
			ring:     &Ring{members: []Node{}},
			node:     node100,
			expected: &Ring{members: []Node{}},
		},

		{
			ring:     &Ring{members: []Node{node100}},
			node:     node100,
			expected: &Ring{members: []Node{}},
		},

		{
			ring:     &Ring{members: []Node{node100}},
			node:     node200,
			expected: &Ring{members: []Node{node100}},
		},

		{
			ring:     &Ring{members: []Node{node100, node200}},
			node:     node200,
			expected: &Ring{members: []Node{node100}},
		},

		{
			ring:     &Ring{members: []Node{node100, node300}},
			node:     node100,
			expected: &Ring{members: []Node{node300}},
		},
	}

	for i, tc := range testCases {
		err := tc.ring.DelNode(tc.node)

		if err != tc.err {
			t.Errorf("%d: got error %v, wanted %v", i, err, tc.err)
		}
		if tc.err != nil {
			continue
		}

		if !reflect.DeepEqual(tc.ring.members, tc.expected.members) {
			t.Errorf("%d: got %+v, wanted %+v", i, tc.ring.members, tc.expected.members)
		}
	}
}
