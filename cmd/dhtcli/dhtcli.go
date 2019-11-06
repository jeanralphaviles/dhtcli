package main

import (
	"github.com/jeanralphaviles/dhtcli/internal/actions"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "dhtcli"
	app.Usage = "Query and interact with the BitTorrent Distributed Hash Table."
	app.Version = "0.0.2"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "bootstrap, b",
			Value: "router.utorrent.com:6881",
			Usage: "Bootstrap DHT node",
		},
	}
	app.Commands = []cli.Command{
		cli.Command{
			Name:      "ping",
			Usage:     "issue a DHT 'ping' to the given node",
			ArgsUsage: "host:port",
			Action:    actions.Ping,
		},
		cli.Command{
			Name:      "find_node",
			Usage:     "issue a DHT 'find_node' request to the given node",
			ArgsUsage: "host:port node_id",
			Action:    actions.FindNode,
		},
		cli.Command{
			Name:      "get_peers",
			Usage:     "issue a DHT 'get_peers' request to the given node",
			ArgsUsage: "host:port info_hash",
			Action:    actions.GetPeers,
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
