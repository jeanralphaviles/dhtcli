package main

import (
	"flag"
	"fmt"
	"github.com/jeanralphaviles/dhtcli/pkg/dht"
	"log"
	"net"
)

var (
	bootstrap = flag.String("bootstrap", "router.utorrent.com:6881", "Bootstrap DHT node")
	infoHash  = flag.String("info_hash", "E2467CBF021192C241367B892230DC1E05C0580E", "Torrent info_hash to query")
)

func main() {
	flag.Parse()
	d, err := dht.New()
	if err != nil {
		log.Fatal(err)
	}
	server, err := net.ResolveUDPAddr("udp", *bootstrap)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := d.Ping(*server)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Ping: %+v\n", resp)
	resp, err = d.GetPeers(*server, *infoHash)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("GetPeers(%q): %+v\n", *infoHash, resp)
	target, ok := resp.Response["id"]
	if !ok {
		log.Fatal(`resp.Response["id"] not found`)
	}
	resp, err = d.FindNode(*server, target)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("FindNode(%q): %+v\n", target, resp)
}
