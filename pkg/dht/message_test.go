package dht

import (
	"encoding/hex"
	"net"
	"reflect"
	"testing"
)

func TestNewRequest(t *testing.T) {
	args := map[string]interface{}{
		"id": "123abc",
	}
	got, err := NewRequest(ping, args)
	if err != nil {
		t.Fatalf("error creating new request: %v", err)
	}
	want := &Message{
		Mtype:     "q",
		Query:     "ping",
		Arguments: args,
		Version:   Version,
	}
	// Set transaction ids to be equal
	want.TransactionID = got.TransactionID
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestNewResponse(t *testing.T) {
	response := map[string]interface{}{
		"id": "123abc",
	}
	got := NewResponse("1234", response)
	want := &Message{
		TransactionID: "1234",
		Mtype:         "r",
		Response:      response,
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestNodes(t *testing.T) {
	cases := []struct {
		in   *Message
		want []Node
		err  bool
	}{
		{
			&Message{Arguments: map[string]interface{}{"nodes": "C4D4E4F5055354C0A801E*E*ii"}},
			[]Node{{ID: []byte("C4D4E4F5055354C0A801"), Peer: &Peer{net.UDPAddr{IP: net.ParseIP("69.42.69.42"), Port: 26985}}}},
			false,
		},
		{
			&Message{Arguments: map[string]interface{}{"nodes": ""}, Response: map[string]interface{}{"nodes": ""}},
			nil,
			true,
		},
		{
			&Message{Response: map[string]interface{}{"nodes": "12345"}},
			nil,
			true,
		},
	}
	for n, c := range cases {
		got, err := c.in.Nodes()
		if (err != nil) != c.err {
			t.Errorf("case %d: expected Nodes() to return error %v, got %v", n, c.err, err)
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("case %d: expected Nodes() = %v, got %v", n, c.want, got)
		}
	}
}

func TestValues(t *testing.T) {
	cases := []struct {
		in   *Message
		want []Peer
		err  bool
	}{
		{
			&Message{Arguments: map[string]interface{}{"values": "*E*Eii"}},
			[]Peer{{net.UDPAddr{IP: net.ParseIP("42.69.42.69"), Port: 26985}}},
			false,
		},
		{
			&Message{Arguments: map[string]interface{}{"values": ""}, Response: map[string]interface{}{"values": ""}},
			nil,
			true,
		},
		{
			&Message{Response: map[string]interface{}{"values": "12345"}},
			nil,
			true,
		},
	}
	for n, c := range cases {
		got, err := c.in.Values()
		if (err != nil) != c.err {
			t.Errorf("case %d: expected Values() to return error %v, got %v", n, c.err, err)
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("case %d: expected Values() = %v, got %v", n, c.want, got)
		}
	}
}

func TestString(t *testing.T) {
	cases := []struct {
		got  *Message
		want string
	}{
		{
			&Message{
				TransactionID: "123",
				Mtype:         "x",
				Query:         string(findNode),
				Arguments: map[string]interface{}{
					"id":     "abc",
					"target": "456",
					"token":  "789",
					"nodes":  "4142434445464748494A4B4C4D4E4F5051525354C0A801010016",
					"values": []interface{}{
						string([]byte{0x7F, 0x00, 0x00, 0x01, 0x00, 0x16}),
					},
				},
				Version: Version,
			}, `{
  "t": "0x313233",
  "y": "x",
  "q": "find_node",
  "a": {
    "id": "0x616263",
    "nodes": [
      {
        "id": "0x3431343234333434343534363437343834393441",
        "address": "52.66.52.67:13380"
      },
      {
        "id": "0x3445344635303531353235333534433041383031",
        "address": "48.49.48.48:12598"
      }
    ],
    "target": "0x343536",
    "token": "0x373839",
    "values": [
      "127.0.0.1:22"
    ]
  },
  "v": "0x2d4443303030312d"
}`}, {
			&Message{
				TransactionID: "123",
				Mtype:         "x",
				Query:         string(findNode),
				Response: map[string]interface{}{
					"id":     "abc",
					"target": "456",
					"token":  "789",
					"nodes":  "12345",
					"values": []interface{}{
						string([]byte{0x7F}),
					},
				},
				Version: Version,
			}, `{
  "t": "0x313233",
  "y": "x",
  "q": "find_node",
  "r": {
    "id": "0x616263",
    "nodes": null,
    "target": "0x343536",
    "token": "0x373839",
    "values": null
  },
  "v": "0x2d4443303030312d"
}`,
		},
	}
	for n, c := range cases {
		if c.want != c.got.String() {
			t.Errorf("case %d: expected %v, got %v", n, c.want, c.got.String())
		}
	}
}

func TestParseCompactNodesEncoding(t *testing.T) {
	encoding, err := hex.DecodeString(
		"4142434445464748494A4B4C4D4E4F5051525354C0A801010016" +
			"6162636465666768696a6b6c6d6e6f70717273747F0000010017",
	)
	if err != nil {
		t.Fatalf("error decoding hex string: %v", err)
	}
	got, err := parseCompactNodesEncoding(encoding)
	if err != nil {
		t.Fatalf("error parsing compact nodes encoding: %v", err)
	}
	want := []Node{{
		ID: []byte("ABCDEFGHIJKLMNOPQRST"),
		Peer: &Peer{
			net.UDPAddr{
				IP:   net.ParseIP("192.168.1.1"),
				Port: 22,
			},
		},
	}, {
		ID: []byte("abcdefghijklmnopqrst"),
		Peer: &Peer{
			net.UDPAddr{
				IP:   net.ParseIP("127.0.0.1"),
				Port: 23,
			},
		},
	}}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected %v, got %v", want, got)
	}

	encoding = []byte("12345")
	if _, err := parseCompactNodesEncoding(encoding); err == nil {
		t.Errorf("parseCompactNodesEncoding(%v) expected error", encoding)
	}
}

func TestParseCompactPeersEncoding(t *testing.T) {
	encoding, err := hex.DecodeString("C0A801010016")
	if err != nil {
		t.Fatalf("error decoding hex string: %v", err)
	}
	got, err := parseCompactPeersEncoding([]interface{}{string(encoding)})
	if err != nil {
		t.Fatalf("error parsing compact peers encoding: %v", err)
	}
	want := []Peer{{
		net.UDPAddr{
			IP:   net.ParseIP("192.168.1.1"),
			Port: 22,
		},
	},
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected %v, got %v", want, got)
	}

	encoding = []byte("12345")
	if _, err := parseCompactPeersEncoding([]interface{}{string(encoding)}); err == nil {
		t.Errorf("parseCompactPeersEncoding(%v) expected error", []interface{}{string(encoding)})
	}
}

func TestParseCompactPeerEncoding(t *testing.T) {
	encoding, err := hex.DecodeString("7F0000010017")
	if err != nil {
		t.Fatalf("error decoding hex string: %v", err)
	}
	got, err := parseCompactPeerEncoding(encoding)
	if err != nil {
		t.Fatalf("error parsing compact peer encoding: %v", err)
	}
	want := &Peer{
		net.UDPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 23,
		},
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected %v, got %v", want, got)
	}

	encoding = []byte("12345")
	if _, err := parseCompactPeerEncoding(encoding); err == nil {
		t.Errorf("parseCompactPeerEncoding(%v) expected error", encoding)
	}
}
