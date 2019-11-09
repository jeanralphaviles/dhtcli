package dht

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/getlantern/deepcopy"
)

const (
	// Version is the arbitrary DHT client string representing dhtcli.
	// Follows the "Peer ID Convention" set forth in BEP 20 with a previously unused client implementation string.
	// https://www.bittorrent.org/beps/bep_0020.html
	Version = "-DC0001-"
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
	TransactionID string                 `bencode:"t" json:"t"`
	Mtype         string                 `bencode:"y" json:"y"`
	Query         string                 `bencode:"q,omitempty" json:"q,omitempty"`
	Arguments     map[string]interface{} `bencode:"a,omitempty" json:"a,omitempty"`
	Response      map[string]interface{} `bencode:"r,omitempty" json:"r,omitempty"`
	Error         []interface{}          `bencode:"e,omitempty" json:"e,omitempty"`
	Version       string                 `bencode:"v,omitempty" json:"v,omitempty"`
}

// newRequest returns a new query message with the specified arguments.
func newRequest(q query, args map[string]interface{}) (*Message, error) {
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
		Version:       Version,
	}, nil
}

// newResponse returns a new response message with the specified response dictionary.
func newResponse(id string, response map[string]interface{}) *Message {
	return &Message{
		TransactionID: id,
		Mtype:         "r",
		Response:      response,
	}
}

// String pretty prints a message as JSON.
func (m *Message) String() string {
	c := &Message{}
	if err := deepcopy.Copy(c, m); err != nil {
		log.Printf("error copying %#v", m)
	}
	// Translate byte strings to hex for better readability.
	c.TransactionID = fmt.Sprintf("0x%x", c.TransactionID)
	c.Version = fmt.Sprintf("0x%x", c.Version)
	for _, dict := range []map[string]interface{}{c.Arguments, c.Response} {
		for k, v := range dict {
			switch k {
			case "nodes":
				n, err := parseCompactNodesEncoding([]byte(v.(string)))
				if err != nil {
					log.Print(err)
				}
				dict["nodes"] = n
			case "values":
				p, err := parseCompactPeersEncoding(v.([]interface{}))
				if err != nil {
					log.Print(err)
				}
				dict["values"] = p
			default:
				dict[k] = fmt.Sprintf("0x%x", v)
			}
		}
	}

	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Printf("error marshalling message: %v", err)
	}
	return string(b)
}

// node encapsulates entries in the "nodes" key in a "find_node" and "get_peers" messages.
type node struct {
	id   []byte
	peer *peer
}

// parseCompactNodesEncoding parses contact information for nodes.
func parseCompactNodesEncoding(b []byte) ([]node, error) {
	buf := bytes.NewBuffer(b)
	if buf.Len()%26 != 0 {
		return nil, fmt.Errorf("compact encoding must be a multiple of 26 bytes long")
	}
	var nodes []node
	for i := 0; i <= buf.Len()/26; i++ {
		id := buf.Next(20)
		peer, err := parseCompactPeerEncoding(buf.Next(6))
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node{
			id:   id,
			peer: peer,
		})
	}
	return nodes, nil
}

func (n *node) MarshalJSON() ([]byte, error) {
	hash := fmt.Sprintf("0x%x", n.id)
	return json.Marshal(
		struct {
			ID   string `json:"id"`
			Peer *peer  `json:"address"`
		}{
			hash,
			n.peer,
		})
}

// peer encapsulates "peer" information included in "find_node" and "get_peers" messages.
type peer struct {
	net.UDPAddr
}

func (p *peer) MarshalText() ([]byte, error) {
	return []byte(p.UDPAddr.String()), nil
}

// parseCompactPeersEncoding parses contact information for peers.
func parseCompactPeersEncoding(e []interface{}) ([]peer, error) {
	var peers []peer
	for _, c := range e {
		peer, err := parseCompactPeerEncoding([]byte(c.(string)))
		if err != nil {
			return nil, err
		}
		peers = append(peers, *peer)
	}
	return peers, nil
}

// parseCompactPeerEncoding parses contact information for a single peer.
func parseCompactPeerEncoding(b []byte) (*peer, error) {
	if len(b) != 6 {
		return nil, fmt.Errorf("compact peer encoding must be 6 bytes long")
	}
	ip := net.IPv4(b[0], b[1], b[2], b[3])
	port := binary.BigEndian.Uint16(b[4:])
	return &peer{
		net.UDPAddr{
			IP:   ip,
			Port: int(port),
		},
	}, nil
}
