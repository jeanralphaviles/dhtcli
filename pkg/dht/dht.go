// Package dht provides methods to query the BitTorrent DHT as defined in BEP 5.
// https://www.bittorrent.org/beps/bep_0005.html
package dht

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"time"

	"github.com/zeebo/bencode"
)

// DHT encapsulates a node in the DHT.
type DHT struct {
	// DHT node id
	id string
}

// New returns a DHT initialized with a random node id.
func New() (*DHT, error) {
	// Generate 20 byte node id.
	id := make([]byte, 20)
	if _, err := rand.Read(id); err != nil {
		return nil, err
	}
	return &DHT{id: string(id)}, nil
}

// query issues a request to a DHT node and returns its response.
func (d *DHT) query(server net.UDPAddr, req *Message) (*Message, error) {
	conn, err := net.DialUDP("udp", nil, &server)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	buf := bytes.NewBuffer([]byte{})
	if err := bencode.NewEncoder(buf).Encode(req); err != nil {
		return nil, fmt.Errorf("error encoding %#v: %v", req, err)
	}
	if _, err := conn.Write(buf.Bytes()); err != nil {
		return nil, err
	}
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	resp := &Message{}
	if err := bencode.NewDecoder(conn).Decode(resp); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}
	return resp, nil
}

// Ping issues a "ping" query to a DHT node and returns its response.
func (d *DHT) Ping(server net.UDPAddr) (*Message, error) {
	args := map[string]interface{}{"id": d.id}
	req, err := newRequest(ping, args)
	if err != nil {
		return nil, fmt.Errorf("error creating ping request: %v", err)
	}
	return d.query(server, req)
}

// encodeInfoHash encodes a string of hexadecimal characters as a string of the literal bytes it represents.
func encodeInfoHash(infoHash string) (string, error) {
	h, err := hex.DecodeString(infoHash)
	if err != nil {
		return "", err
	}
	if len(h) != 20 {
		return "", fmt.Errorf("invalid infoHash %q: needs to represent 20 bytes", infoHash)
	}
	return string(h), nil
}

// FindNode issues a "find_node" query to a DHT node and returns its response.
func (d *DHT) FindNode(server net.UDPAddr, target string) (*Message, error) {
	hash, err := encodeInfoHash(target)
	if err != nil {
		return nil, err
	}
	args := map[string]interface{}{"id": d.id, "target": hash}
	req, err := newRequest(findNode, args)
	if err != nil {
		return nil, fmt.Errorf("error creating find_node request: %v", err)
	}
	return d.query(server, req)
}

// GetPeers issues a "get_peers" query to a DHT node and returns its response.
func (d *DHT) GetPeers(server net.UDPAddr, infoHash string) (*Message, error) {
	infoHash, err := encodeInfoHash(infoHash)
	if err != nil {
		return nil, err
	}
	args := map[string]interface{}{"id": d.id, "info_hash": infoHash}
	req, err := newRequest(getPeers, args)
	if err != nil {
		return nil, fmt.Errorf("error creating get_peers request: %v", err)
	}
	return d.query(server, req)
}

// encodeToken encodes a string of hexadecimal characters of a token as the literal bytes it represents.
func encodeToken(token string) (string, error) {
	h, err := hex.DecodeString(token)
	if err != nil {
		return "", err
	}
	return string(h), nil
}

// AnnouncePeer issues an "announce_peer" query to a DHT node and returns its response.
//
// infoHash is the hash of the torrent to announce as a peer of.
// token is the token received in a previous get_peers request to this server.
// port is the intended server port of this peer. If zero, the "implied_port" setting will be sent in the request.
func (d *DHT) AnnouncePeer(server net.UDPAddr, infoHash string, token string, port int) (*Message, error) {
	infoHash, err := encodeInfoHash(infoHash)
	if err != nil {
		return nil, fmt.Errorf("error encoding infoHash: %v", err)
	}
	token, err = encodeToken(token)
	if err != nil {
		return nil, fmt.Errorf("error encoding token: %v", err)
	}
	impliedPort := 0
	if port == 0 {
		impliedPort = 1
	}
	args := map[string]interface{}{
		"id":           d.id,
		"implied_port": impliedPort,
		"info_hash":    infoHash,
		"port":         port,
		"token":        token,
	}
	req, err := newRequest(announcePeer, args)
	if err != nil {
		return nil, fmt.Errorf("error creating announce_peer request: %v", err)
	}
	return d.query(server, req)
}
