// Package actions provides handlers for dhtcli commands.
package actions

import (
	"fmt"
	"github.com/jeanralphaviles/dhtcli/pkg/dht"
	"github.com/urfave/cli"
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
