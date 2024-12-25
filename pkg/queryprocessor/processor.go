package queryprocessor

import (
	"bytes"
	"fmt"
	"github.com/jeanralphaviles/dhtcli/pkg/dht"
	"log"
	"math/big"
	"net"
)

// QueryProcessor maintains state for queries into the DHT.
type QueryProcessor struct {
	dht          *dht.DHT
	routingTable *routingTable
}

// New returns a new DHT QueryProcessor initialized with a bootstrap node.
func New(bootstrap net.UDPAddr, k int) (*QueryProcessor, error) {
	d, err := dht.New()
	if err != nil {
		return nil, fmt.Errorf("error creating DHT object: %v", err)
	}
	rt, err := newRoutingTable(k)
	if err != nil {
		return nil, fmt.Errorf("error creating routing table: %v", err)
	}
	q := &QueryProcessor{
		dht:          d,
		routingTable: rt,
	}
	// Get node id of Bootstrap node.
	resp, err := d.Ping(bootstrap)
	if err != nil {
		return nil, fmt.Errorf("error determining id of bootstrap node: %v", err)
	}
	id, ok := resp.Response["id"]
	if !ok {
		return nil, fmt.Errorf("ping response from bootstrap node did not include id: %v", resp)
	}
	node := dht.Node{
		ID: []byte(id.(string)),
		Peer: &dht.Peer{
			UDPAddr: bootstrap,
		},
	}
	distance := big.NewInt(0)
	q.routingTable.insert(node, *distance)
	return q, nil
}

// distance returns the distance metric between two node ids.
//
// distance is defined as the XOR of two ids interpreted as an integer.
// distance(a, b) = |A XOR B|
func distance(a, b []byte) (*big.Int, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("distance must be called with ids of the same length")
	}
	buf := make([]byte, len(a))
	for i := range a {
		buf[i] = a[i] ^ b[i]
	}
	i := new(big.Int)
	i.SetBytes(buf)
	return i, nil
}

// FindNode finds the contact information for a target node given its node id.
//
// Returns a response with a key "nodes" containing information for the target
// node and/or the closest nodes to the target.
func (q *QueryProcessor) FindNode(target string) (*dht.Message, error) {
	var ret *dht.Message
	closestDistance := new(big.Int)
	// Furthest distance possible, 2^160.
	closestDistance.SetBytes(bytes.Repeat([]byte{0xFF}, 20))
	var visited []dht.Node
	for q.routingTable.Len() > 0 {
		node, _ := q.routingTable.pop()
		visited = append(visited, node)
		resp, err := q.dht.FindNode(node.Peer.UDPAddr, target)
		if err != nil {
			log.Print(err)
			continue
		}
		nodes, err := resp.Nodes()
		if err != nil {
			log.Print(err)
			continue
		}
		if ret == nil {
			// At least return the result of the first FindNode query if it
			// contains any nodes.
			ret = resp
		}
		t, _ := dht.EncodeInfoHash(target)
		d, err := distance([]byte(t), node.ID)
		if err != nil {
			log.Printf("distance(%x, %x): %v", []byte(t), node.ID, err)
			continue
		}
		if d.Cmp(closestDistance) < 0 {
			// Set ret to the FindNodes response from the closest node heard from.
			ret = resp
			closestDistance = d
		}
	NODES:
		for _, n := range nodes {
			distance, err := distance([]byte(t), n.ID)
			if err != nil {
				log.Printf("distance(%x, %x): %v", []byte(t), n.ID, err)
				continue
			}
			if distance.Int64() == 0 {
				// Target found
				return ret, nil
			}
			// Exclude previously visited nodes.
			for _, v := range visited {
				if bytes.Equal(v.ID, n.ID) {
					continue NODES
				}
			}
			q.routingTable.insert(n, *distance)
		}
	}
	if ret == nil {
		return nil, fmt.Errorf("could not successfully query any DHT nodes")
	}
	return ret, nil
}
