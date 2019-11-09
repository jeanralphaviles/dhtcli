package queryprocessor

import (
	"github.com/jeanralphaviles/dhtcli/pkg/dht"
	"math/big"
	"sort"
	"testing"
)

func TestNewRoutingTable(t *testing.T) {
	cases := []struct {
		size int
		fail bool
	}{
		{1, false},
		{0, true},
	}
	for n, c := range cases {
		if _, err := newRoutingTable(c.size); (err != nil) != c.fail {
			t.Errorf("case %d: expected newRoutingTable(%d) to return error: %v, got %v", n, c.size, c.fail, err)
		}
	}
}

func TestInsert(t *testing.T) {
	nodes := make([]dht.Node, 3)
	rt, _ := newRoutingTable(2)
	for i, n := range nodes {
		// Decreasing priority to assure sort must run.
		rt.insert(n, *big.NewInt(int64(3 - i)))
	}
	if rt.Len() != 2 {
		t.Errorf("routing table should contain 2 entries, has %d", rt.Len())
	}
	if !sort.IsSorted(rt) {
		t.Errorf("routing table should be sorted")
	}
}

func TestPop(t *testing.T) {
	cases := []struct {
		nodes []dht.Node
		fail  bool
	}{
		{[]dht.Node{dht.Node{}}, false},
		{[]dht.Node{}, true},
	}
	for n, c := range cases {
		rt, _ := newRoutingTable(1)
		for i, n := range c.nodes {
			rt.insert(n, *big.NewInt(int64(i)))
		}
		if _, err := rt.pop(); (err != nil) != c.fail {
			t.Errorf("case %d: expected rt.Pop to return error %v, got %v", n, c.fail, err)
		}
	}
}
