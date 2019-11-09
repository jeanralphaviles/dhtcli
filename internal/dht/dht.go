// Package dht contains handlers for dhtcli dht commands.
package dht

import (
	"fmt"
	"github.com/jeanralphaviles/dhtcli/pkg/queryprocessor"
	"github.com/urfave/cli"
	"net"
)

// FindNode searches the BitTorrent DHT for the contact information of a target node.
func FindNode(c *cli.Context) error {
	if c.NArg() != 1 {
		command := c.Command
		return fmt.Errorf("%v: %v", command.FullName(), command.ArgsUsage)
	}
	bootstrap, err := resolveBootstrap(c)
	if err != nil {
		return err
	}
	q, err := queryprocessor.New(*bootstrap, c.Int("table_size"))
	if err != nil {
		return err
	}
	target := c.Args().Get(0)
	resp, err := q.FindNode(target)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", resp)
	return nil
}

func resolveBootstrap(c *cli.Context) (*net.UDPAddr, error) {
	bootstrap, err := net.ResolveUDPAddr("udp", c.String("bootstrap"))
	if err != nil {
		return nil, fmt.Errorf("error resolving bootstrap node: %v", err)
	}
	return bootstrap, err
}
