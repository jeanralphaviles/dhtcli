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
	app.Version = "0.0.4"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "bootstrap, b",
			Value: "router.utorrent.com:6881",
			Usage: "Bootstrap DHT node",
		},
	}
	app.Commands = []cli.Command{
		cli.Command{
			Name:  "query",
			Usage: "Issue individual requests to a DHT node.",
			Subcommands: []cli.Command{
				cli.Command{
					Name:      "ping",
					Usage:     "Issue a DHT 'ping' to the given node",
					ArgsUsage: "host:port",
					Action:    actions.Ping,
				},
				cli.Command{
					Name:      "find_node",
					Usage:     "Issue a DHT 'find_node' request to the given node",
					ArgsUsage: "host:port node_id",
					Action:    actions.FindNode,
				},
				cli.Command{
					Name:      "get_peers",
					Usage:     "Issue a DHT 'get_peers' request to the given node",
					ArgsUsage: "host:port info_hash",
					Action:    actions.GetPeers,
				},
				cli.Command{
					Name:      "announce_peer",
					Usage:     "Issue a DHT 'announce_peer' request to the given node",
					ArgsUsage: "host:port info_hash",
					Description: "Announce this IP:Port to a DHT node as a peer for the torrent with the given info_hash. \n\n" +
						"   This request requires a token received from the node in a previous get_peers request. " +
						"This token is used to prevent malicious hosts from signing up other hosts for torrents. " +
						"If --token is not specified, one will be obtained by issuing a get_peers request to the node. \n\n" +
						"   --port specifies the port of the announced peer. " +
						"If port is set to 0, the announce_peer request will contain the \"implied_port\" setting. " +
						"This setting will derive the port value automatically as described in BEP 5.",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "token, t",
							Usage: "Response token from previous get_peers request to this node",
						},
						cli.IntFlag{
							Name:  "port, p",
							Usage: "Port value of announced peer",
							Value: 0,
						},
					},
					Action: actions.AnnouncePeer,
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
