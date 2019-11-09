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
	d, err := New()
	if err != nil {
		t.Fatalf("error creating new DHT object: %v", err)
	}
	got, err := d.Ping(*addr)
	if err != nil {
		t.Fatalf("error issuing Ping: %v", err)
	}
	want, err := newRequest(ping, map[string]interface{}{
		"id": d.ID,
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
	d, err := New()
	if err != nil {
		t.Fatalf("error creating new DHT object: %v", err)
	}
	got, err := d.FindNode(*addr, "4142434445464748494A4B4C4D4E4F5051525354")
	if err != nil {
		t.Fatalf("error issuing FindNode: %v", err)
	}
	want, err := newRequest(findNode, map[string]interface{}{
		"id":     d.ID,
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

	errCases := []struct {
		addr   *net.UDPAddr
		target string
	}{
		{addr, "ABCDEFG"},
		{addr, "12"},
		{&net.UDPAddr{Port: 0}, "4142434445464748494A4B4C4D4E4F5051525354"},
	}
	for n, c := range errCases {
		if _, err := d.FindNode(*c.addr, c.target); err == nil {
			t.Errorf("case %d: d.FindNode(%v, %q) did not error", n, c.addr, c.target)
		}
	}
}

func TestGetPeers(t *testing.T) {
	d, err := New()
	if err != nil {
		t.Fatalf("error creating new DHT object: %v", err)
	}
	got, err := d.GetPeers(*addr, "4142434445464748494A4B4C4D4E4F5051525354")
	if err != nil {
		t.Fatalf("error issuing GetPeers: %v", err)
	}
	want, err := newRequest(getPeers, map[string]interface{}{
		"id":        d.ID,
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

	errCases := []struct {
		addr   *net.UDPAddr
		target string
	}{
		{addr, "ABCDEFG"},
		{addr, "12"},
		{&net.UDPAddr{Port: 65536}, "4142434445464748494A4B4C4D4E4F5051525354"},
	}
	for n, c := range errCases {
		if _, err := d.GetPeers(*c.addr, c.target); err == nil {
			t.Errorf("case %d: d.GetPeers(%v, %q) did not error", n, c.addr, c.target)
		}
	}
}

func TestAnnouncePeer(t *testing.T) {
	d, err := New()
	if err != nil {
		t.Fatalf("error creating new DHT object: %v", err)
	}
	cases := []struct {
		port int
	}{
		{0},
		{6881},
	}
	for n, c := range cases {
		got, err := d.AnnouncePeer(*addr, "4142434445464748494A4B4C4D4E4F5051525354", "746F6B656E", c.port)
		if err != nil {
			t.Fatalf("case %d: error issuing AnnouncePeer request: %v", n, err)
		}
		impliedPort := 0
		if c.port == 0 {
			impliedPort = 1
		}
		want, err := newRequest(announcePeer, map[string]interface{}{
			"id":           d.ID,
			"implied_port": int64(impliedPort),
			"info_hash":    "ABCDEFGHIJKLMNOPQRST",
			"port":         int64(c.port),
			"token":        "token",
		})
		// Set transaction ids to be equal
		want.TransactionID = got.TransactionID
		if !reflect.DeepEqual(want, got) {
			t.Errorf("case %d: expected %v, got %v", n, want, got)
		}
	}

	errCases := []struct {
		infoHash string
		token    string
	}{
		{"4142434445464748494A4B4C4D4E4F5051525354", "ABCDEFG"},
		{"ABCDEFG", "746F6B656E"},
	}
	for n, c := range errCases {
		if _, err := d.AnnouncePeer(*addr, c.infoHash, c.token, 0); err == nil {
			t.Errorf("case %d: expected d.AnnouncePeer(%v, %q, %q) to error", n, addr, c.infoHash, c.token)
		}
	}
}
