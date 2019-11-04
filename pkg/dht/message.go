package dht

import (
	"crypto/rand"
)

type query string

// Message query types.
const (
	announcePeer query = "announce_peer"
	findNode     query = "find_node"
	getPeers     query = "get_peers"
	ping         query = "ping"
)

// Message encapsulates a DHT message as defined in BEP 5.
type Message struct {
	TransactionID string            `bencode:"t"`
	Mtype         string            `bencode:"y"`
	Query         string            `bencode:"q,omitempty"`
	Arguments     map[string]string `bencode:"a,omitempty"`
	Response      map[string]string `bencode:"r,omitempty"`
	Error         []interface{}     `bencode:"e,omitempty"`
	Version       string            `bencode:"v,omitempty"`
}

// newRequest returns a new query message with the specified arguments.
func newRequest(q query, args map[string]string) (*Message, error) {
	// 2 bytes recommended by BEP 5.
	id := make([]byte, 2)
	if _, err := rand.Read(id); err != nil {
		return nil, err
	}
	return &Message{
		TransactionID: string(id),
		Mtype:         "q",
		Query:         string(q),
		Arguments:     args,
		Version:       "-DC0001-", // Arbitrary client string representing dhtcli.
	}, nil
}

// newResponse returns a new response message with the specified response
// dictionary.
func newResponse(id string, response map[string]string) *Message {
	return &Message{
		TransactionID: id,
		Mtype:         "r",
		Response:      response,
	}
}
