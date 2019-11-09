package queryprocessor

import (
	"fmt"
	"github.com/jeanralphaviles/dhtcli/pkg/dht"
	"math/big"
	"sort"
)

// routingTable maintains a list of known good nodes used as starting points
// for queries into the DHT.
type routingTable struct {
	// Maximum number of nodes to maintain
	size    int
	entries []entry
}

type entry struct {
	node     dht.Node
	distance big.Int
}

func newRoutingTable(size int) (*routingTable, error) {
	if size <= 0 {
		return nil, fmt.Errorf("routing table size must be >= 1, got %d", size)
	}
	return &routingTable{
		size: size,
	}, nil
}

func (r *routingTable) insert(n dht.Node, distance big.Int) {
	entry := entry{n, distance}
	r.entries = append(r.entries, entry)
	sort.Sort(r)
	// Trim to size.
	if len(r.entries) > r.size {
		r.entries = r.entries[:r.size]
	}
}

func (r *routingTable) pop() (dht.Node, error) {
	if len(r.entries) == 0 {
		return dht.Node{}, fmt.Errorf("pop() called on empty routing table")
	}
	var ret dht.Node
	ret, r.entries = r.entries[0].node, r.entries[1:]
	return ret, nil
}

func (r *routingTable) Len() int {
	return len(r.entries)
}

func (r *routingTable) Less(i, j int) bool {
	return r.entries[i].distance.Cmp(&r.entries[j].distance) < 0
}

func (r *routingTable) Swap(i, j int) {
	r.entries[i], r.entries[j] = r.entries[j], r.entries[i]
}
