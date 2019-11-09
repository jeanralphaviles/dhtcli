package queryprocessor

import (
	"bytes"
	"github.com/jeanralphaviles/dhtcli/pkg/dht"
	"github.com/zeebo/bencode"
	"log"
	"math/big"
	"net"
	"reflect"
	"testing"
)

var addr *net.UDPAddr
var buffer *bytes.Buffer

func init() {
	var err error
	addr, err = net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		log.Panic(err)
	}
	server, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Panic(err)
	}
	addr, err = net.ResolveUDPAddr("udp", server.LocalAddr().String())
	if err != nil {
		log.Panic(err)
	}
	buffer = bytes.NewBuffer([]byte{})
	go func() {
		defer server.Close()
		buf := make([]byte, 1024)
		for {
			_, client, err := server.ReadFromUDP(buf)
			if err != nil {
				log.Panic(err)
			}
			_, err = server.WriteTo(buffer.Bytes(), client)
			if err != nil {
				log.Panic(err)
			}
			buffer.Reset()
		}
	}()
}

func TestNew(t *testing.T) {
	cases := []struct {
		resp map[string]interface{}
		size int
		fail bool
	}{
		{map[string]interface{}{"id": "1234"}, 1, false},
		{map[string]interface{}{"not id": ""}, 1, true},
		{nil, 0, true},
	}
	for n, c := range cases {
		m := dht.NewResponse("123", c.resp)
		if err := bencode.NewEncoder(buffer).Encode(m); err != nil {
			t.Errorf("case %d: error writing to buffer: %v", n, err)
			continue
		}
		_, err := New(*addr, c.size)
		if (err != nil) != c.fail {
			t.Errorf("case %d: expected New() to return error: %v, got %v", n, c.fail, err)
			continue
		}
		if c.fail {
			buffer.Reset()
		}
	}
}

func TestDistance(t *testing.T) {
	cases := []struct {
		a    []byte
		b    []byte
		want *big.Int
		fail bool
	}{
		{[]byte{0xFF, 0xFF}, []byte{0xFF, 0xBA}, big.NewInt(69), false},
		{[]byte{0xFF, 0xFF}, []byte{0xFF}, nil, true},
	}
	for n, c := range cases {
		got, err := distance(c.a, c.b)
		if (err != nil) != c.fail {
			t.Errorf("case %d: expected distance(0x%x, 0x%x) to return error: %v, got %v", n, c.a, c.b, c.fail, err)
			continue
		}
		if c.fail {
			continue
		}
		if got.Cmp(c.want) != 0 {
			t.Errorf("case %d: distance(0x%x, 0x%x) = %d, want %d", n, c.a, c.b, got, c.want)
			continue
		}
	}
}

func TestFindNode(t *testing.T) {
	d, err := dht.New()
	if err != nil {
		t.Fatalf("error calling dht.New(): %v", err)
	}
	q := &QueryProcessor{
		dht: d,
	}
	cases := []struct {
		bootstrap *dht.Node
		target    string
		resp      map[string]interface{}
		want      map[string]interface{}
		fail      bool
	}{
		{
			// FindNode called once, it's returned.
			bootstrap: &dht.Node{Peer: &dht.Peer{*addr}},
			target:    "4142434445464748494A4B4C4D4E4F5051525354",
			resp:      map[string]interface{}{"nodes": "C4D4E4F5055354C0A801E*E*ii"},
			want:      map[string]interface{}{"nodes": "C4D4E4F5055354C0A801E*E*ii"},
			fail:      false,
		},
		{
			// Already visited, response.
			bootstrap: &dht.Node{ID: []byte("C4D4E4F5055354C0A801"), Peer: &dht.Peer{*addr}},
			target:    "4142434445464748494A4B4C4D4E4F5051525354",
			resp:      map[string]interface{}{"nodes": "C4D4E4F5055354C0A801E*E*ii"},
			want:      map[string]interface{}{"nodes": "C4D4E4F5055354C0A801E*E*ii"},
			fail:      false,
		},
		{
			// Target found
			bootstrap: &dht.Node{ID: []byte("C4D4E4F5055354C0A801"), Peer: &dht.Peer{*addr}},
			target:    "4142434445464748494A4B4C4D4E4F5051525354",
			resp:      map[string]interface{}{"nodes": "ABCDEFGHIJKLMNOPQRSTE*E*ii"},
			want:      map[string]interface{}{"nodes": "ABCDEFGHIJKLMNOPQRSTE*E*ii"},
			fail:      false,
		},
		{
			// Target is incorrect length.
			bootstrap: &dht.Node{ID: []byte("C4D4E4F5055354C0A801"), Peer: &dht.Peer{*addr}},
			target:    "123",
			fail:      true,
		},
		{
			// Returned node has incorrect length.
			bootstrap: &dht.Node{Peer: &dht.Peer{*addr}},
			target:    "4142434445464748494A4B4C4D4E4F5051525354",
			resp:      map[string]interface{}{"nodes": "1234"},
			fail:      true,
		},
		{
			// Routing table is empty.
			bootstrap: nil,
			fail:      true,
		},
	}
	for n, c := range cases {
		rt, err := newRoutingTable(1)
		if err != nil {
			t.Fatalf("error creating routing table: %v", err)
		}
		if c.bootstrap != nil {
			rt.insert(*c.bootstrap, *big.NewInt(1000))
		}
		q.routingTable = rt
		m := dht.NewResponse("", c.resp)
		if err := bencode.NewEncoder(buffer).Encode(m); err != nil {
			t.Errorf("case %d: error writing to buffer: %v", n, err)
			continue
		}
		got, err := q.FindNode(c.target)
		if (err != nil) != c.fail {
			t.Errorf("case %d: expected q.FindNode(%q) to return error: %v, got %v", n, c.target, c.fail, err)
			continue
		}
		if c.fail {
			continue
		}
		want := dht.NewResponse("", c.want)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("case %d: got %v, want %v", n, got, want)
		}
	}
}
