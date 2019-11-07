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
	got, err := newRequest(ping, args)
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
	got := newResponse("1234", response)
	want := &Message{
		TransactionID: "1234",
		Mtype:         "r",
		Response:      response,
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected %v, got %v", want, got)
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
	want := []node{{
		id: []byte("ABCDEFGHIJKLMNOPQRST"),
		peer: &peer{
			net.UDPAddr{
				IP:   net.ParseIP("192.168.1.1"),
				Port: 22,
			},
		},
	}, {
		id: []byte("abcdefghijklmnopqrst"),
		peer: &peer{
			net.UDPAddr{
				IP:   net.ParseIP("127.0.0.1"),
				Port: 23,
			},
		},
	}}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected %v, got %v", want, got)
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
	want := []peer{{
		net.UDPAddr{
			IP:   net.ParseIP("192.168.1.1"),
			Port: 22,
		},
	},
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected %v, got %v", want, got)
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
	want := &peer{
		net.UDPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 23,
		},
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected %v, got %v", want, got)
	}
}
