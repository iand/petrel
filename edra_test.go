package petrel

import (
	"net"
	"reflect"
	"testing"
)

func TestAddMember(t *testing.T) {

	addr100 := address{id: 100, addr: &net.IPAddr{IP: net.IPv4(1, 1, 1, 1)}}
	addr200 := address{id: 200, addr: &net.IPAddr{IP: net.IPv4(1, 1, 1, 2)}}
	addr300 := address{id: 300, addr: &net.IPAddr{IP: net.IPv4(1, 1, 1, 3)}}

	testCases := []struct {
		before   []address
		addr     address
		expected []address
	}{
		{
			before:   []address{},
			addr:     addr100,
			expected: []address{addr100},
		},

		{
			before:   []address{addr100},
			addr:     addr200,
			expected: []address{addr100, addr200},
		},

		{
			before:   []address{addr100},
			addr:     addr100,
			expected: []address{addr100},
		},

		{
			before:   []address{addr200},
			addr:     addr100,
			expected: []address{addr100, addr200},
		},

		{
			before:   []address{addr100, addr300},
			addr:     addr200,
			expected: []address{addr100, addr200, addr300},
		},
	}

	for i, tc := range testCases {
		in := instance{
			memberList: tc.before,
		}
		in.addMember(tc.addr)

		if !reflect.DeepEqual(in.memberList, tc.expected) {
			t.Errorf("%d: got %+v, wanted %+v", i, in.memberList, tc.expected)
		}
	}
}

func TestDelMember(t *testing.T) {

	addr100 := address{id: 100, addr: &net.IPAddr{IP: net.IPv4(1, 1, 1, 1)}}
	addr200 := address{id: 200, addr: &net.IPAddr{IP: net.IPv4(1, 1, 1, 2)}}
	addr300 := address{id: 300, addr: &net.IPAddr{IP: net.IPv4(1, 1, 1, 3)}}

	testCases := []struct {
		before   []address
		addr     address
		expected []address
	}{
		{
			before:   []address{},
			addr:     addr100,
			expected: []address{},
		},

		{
			before:   []address{addr100},
			addr:     addr100,
			expected: []address{},
		},

		{
			before:   []address{addr100},
			addr:     addr200,
			expected: []address{addr100},
		},

		{
			before:   []address{addr100, addr200},
			addr:     addr200,
			expected: []address{addr100},
		},

		{
			before:   []address{addr100, addr300},
			addr:     addr100,
			expected: []address{addr300},
		},
	}

	for i, tc := range testCases {
		in := instance{
			memberList: tc.before,
		}
		in.delMember(tc.addr)

		if !reflect.DeepEqual(in.memberList, tc.expected) {
			t.Errorf("%d: got %+v, wanted %+v", i, in.memberList, tc.expected)
		}
	}
}

func TestLog2Ceil(t *testing.T) {
	testCases := []struct {
		n int
		x int
	}{
		{0, 0},
		{1, 0},
		{2, 1},
		{3, 2},
		{4, 2},
		{5, 3},
		{6, 3},
		{7, 3},
		{8, 3},
		{9, 4},
		{31, 5},
		{32, 5},
		{33, 6},
		{453, 9},
	}

	for _, tc := range testCases {
		x := log2ceil(tc.n)
		if x != tc.x {
			t.Errorf("%d: got %d, wanted %d", tc.n, x, tc.x)
		}
	}

}

func TestPow2(t *testing.T) {
	testCases := []struct {
		n int
		x int
	}{
		{0, 1},
		{1, 2},
		{2, 4},
		{3, 8},
		{4, 16},
		{5, 32},
		{6, 64},
		{7, 128},
		{8, 256},
	}

	for _, tc := range testCases {
		x := pow2(tc.n)
		if x != tc.x {
			t.Errorf("%d: got %d, wanted %d", tc.n, x, tc.x)
		}
	}

}
