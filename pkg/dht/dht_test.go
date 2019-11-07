package dht

import (
	"log"
	"net"
	"reflect"
	"testing"
)

var addr *net.UDPAddr

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
	go func() {
		defer server.Close()
		buf := make([]byte, 1024)
		for {
			n, client, err := server.ReadFromUDP(buf)
			if err != nil {
				log.Panic(err)
			}
			_, err = server.WriteTo(buf[:n], client)
			if err != nil {
				log.Panic(err)
			}
		}
	}()
}

func TestPing(t *testing.T) {
	d := DHT{id: "abc123"}
	got, err := d.Ping(*addr)
	if err != nil {
		t.Fatalf("error issuing Ping: %v", err)
	}
	want, err := newRequest(ping, map[string]interface{}{
		"id": "abc123",
	})
	if err != nil {
		t.Fatalf("error creating Ping request: %v", err)
	}
	// Set transaction ids to be equal
	want.TransactionID = got.TransactionID
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestFindNode(t *testing.T) {
	d := DHT{id: "abc123"}
	got, err := d.FindNode(*addr, "4142434445464748494A4B4C4D4E4F5051525354")
	if err != nil {
		t.Fatalf("error issuing FindNode: %v", err)
	}
	want, err := newRequest(findNode, map[string]interface{}{
		"id":     "abc123",
		"target": "ABCDEFGHIJKLMNOPQRST",
	})
	if err != nil {
		t.Fatalf("error creating FindNode request: %v", err)
	}
	// Set transaction ids to be equal
	want.TransactionID = got.TransactionID
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestGetPeers(t *testing.T) {
	d := DHT{id: "abc123"}
	got, err := d.GetPeers(*addr, "4142434445464748494A4B4C4D4E4F5051525354")
	if err != nil {
		t.Fatalf("error issuing GetPeers: %v", err)
	}
	want, err := newRequest(getPeers, map[string]interface{}{
		"id":        "abc123",
		"info_hash": "ABCDEFGHIJKLMNOPQRST",
	})
	if err != nil {
		t.Fatalf("error creating GetPeers request: %v", err)
	}
	// Set transaction ids to be equal
	want.TransactionID = got.TransactionID
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected %v, got %v", want, got)
	}
}
