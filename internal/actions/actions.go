// Package actions provides handlers for dhtcli commands.
package actions

import (
	"fmt"
	"github.com/jeanralphaviles/dhtcli/pkg/dht"
	"github.com/urfave/cli"
	"log"
	"net"
)

// Ping issues a "ping" query to a DHT node and prints its response.
func Ping(c *cli.Context) error {
	if c.NArg() != 1 {
		command := c.Command
		return fmt.Errorf("%v: %v", command.FullName(), c.Command.ArgsUsage)
	}
	server, err := net.ResolveUDPAddr("udp", c.Args().Get(0))
	if err != nil {
		return err
	}
	d, err := dht.New()
	if err != nil {
		return err
	}
	resp, err := d.Ping(*server)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", resp)
	return nil
}

// FindNode issues a "find_node" query to a DHT node and prints its response.
func FindNode(c *cli.Context) error {
	if c.NArg() != 2 {
		command := c.Command
		return fmt.Errorf("%v: %v", command.FullName(), c.Command.ArgsUsage)
	}
	server, err := net.ResolveUDPAddr("udp", c.Args().Get(0))
	if err != nil {
		return err
	}
	d, err := dht.New()
	if err != nil {
		return err
	}
	resp, err := d.FindNode(*server, c.Args().Get(1))
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", resp)
	return nil
}

// GetPeers issues a "get_peers" query to a DHT node and prints its response.
func GetPeers(c *cli.Context) error {
	if c.NArg() != 2 {
		command := c.Command
		return fmt.Errorf("%v: %v", command.FullName(), c.Command.ArgsUsage)
	}
	server, err := net.ResolveUDPAddr("udp", c.Args().Get(0))
	if err != nil {
		return err
	}
	d, err := dht.New()
	if err != nil {
		return err
	}
	resp, err := d.GetPeers(*server, c.Args().Get(1))
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", resp)
	return nil
}

// AnnouncePeer issues an "announce_peer" request to a DHT node and prints its response.
//
// infoHash contains the 20-byte hash of the target torrent and token is the
// token received from a previous "get_peers" request. If a token is not
// specified, a GetPeers request will be issued to obtain one.
func AnnouncePeer(c *cli.Context) error {
	if c.NArg() != 2 {
		command := c.Command
		return fmt.Errorf("%v: %v", command.FullName(), c.Command.ArgsUsage)
	}
	server, err := net.ResolveUDPAddr("udp", c.Args().Get(0))
	if err != nil {
		return err
	}
	hash := c.Args().Get(1)
	d, err := dht.New()
	if err != nil {
		return err
	}
	token := c.String("token")
	if token == "" {
		log.Print("--token not specified, issuing get_peers request first to obtain one.")
		resp, err := d.GetPeers(*server, hash)
		if err != nil {
			return err
		}
		var ok bool
		token, ok = resp.Response["token"].(string)
		if !ok {
			return fmt.Errorf("token not present in response: %v", resp)
		}
		// d.AnnouncePeer expects token as a hex string.
		token = fmt.Sprintf("%x", token)
		log.Printf("Got token 0x%v.", token)
	}
	port := c.Int("port")
	resp, err := d.AnnouncePeer(*server, hash, token, port)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", resp)
	return nil
}
