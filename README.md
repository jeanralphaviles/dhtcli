# dhtcli - BitTorrent Distributed Hash Table CLI

Query and interact with the BitTorrent DHT with a simple CLI tool. Implements
the [BEP 5](https://www.bittorrent.org/beps/bep_0005.html) DHT protocol
specification.

Returns results as JSON for easy parsing.

## Installation

You need to have a working Go environment with version 1.13 or greater
[installed](https://golang.org/doc/install).

```shell
GO111MODULE=on go get -u github.com/jeanralphaviles/dhtcli/...
```

## Usage

```bash
NAME:
   dhtcli - Query and interact with the BitTorrent Distributed Hash Table.

USAGE:
   dhtcli [global options] command [command options] [arguments...]

VERSION:
   0.0.5

COMMANDS:
   query    Issue individual requests to a BitTorrent DHT node.
   dht      [Experimental] - Issues requests to the BitTorrent DHT.
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### Example

```shell
$ dhtcli query get_peers router.utorrent.com:6881 E2467CBF021192C241367B892230DC1E05C0580E
{
  "t": "0x0d4e",
  "y": "r",
  "r": {
    "id": "0xe246e7edc7cf68b7f3abf20ddb042f8182312e0e",
    "nodes": [
      {
        "id": "0xe2462bdd340f0ae1b51a7b1d98f3ad8e6e33e7ac",
        "address": "189.101.214.57:13358"
      },
      {
        "id": "0xe24678d6ae529049f1f1bbe9ebb3a6db3c870ce1",
        "address": "85.66.198.68:5794"
      },
      {
        "id": "0xe2460f49f1f1bbe9ebb3a6db3c870c3e99245e52",
        "address": "95.149.5.53:42061"
      },
      {
        "id": "0xe246555b4be3cf1d716c845218e9fe54feaa9c65",
        "address": "101.165.148.218:18437"
      },
      {
        "id": "0xe246405e58e82ba515172f8d0e2e69d00ea88a63",
        "address": "83.16.219.194:6889"
      }
    ],
    "p": "0x961f",
    "token": "0x29e67ee7",
    "values": [
      "39.8.43.112:25080",
      "41.83.3.125:23227",
      "47.213.57.140:23985",
      "59.46.188.197:19103",
      "68.133.23.82:52481",
      "78.56.91.69:49369",
      "86.210.223.121:31316",
      "90.146.77.100:9225",
      "94.212.49.150:59831",
      "103.106.150.134:43895",
      "107.15.138.17:37120",
      "166.173.62.88:44149",
      "179.186.148.131:9668",
      "197.148.1.17:49291"
    ]
  },
  "v": "0x4c540100"
}
```

Note, this won't always get you peers. If the queried node has peers for the
info_hash, they are returned in a key "values" as a list of strings. Values are
in IP:Port format. If the queried node has no peers for the info_hash, a key
"nodes" is returned containing the K closest nodes to the info_hash.

Issue subsequent get_peers requests to the returned nodes until peers are
found.

## Subcommands

### Query

The query subcommands issues individual requests to a BitTorrent DHT node.

#### ping

The most basic query is ping.

A server should respond with a single key "id", containing the queried node's
ID.

```shell
$ dhtcli query ping router.utorrent.com:6881
{
  "t": "0xc4be",
  "y": "r",
  "r": {
    "id": "0xebff36697351ff4aec29cdbaabf2fbe3467cc267"
  },
  "v": "0x"
}
```

#### find_node

Find node is used to find the contact information for a node given its ID.

Note: an node ID may also be an info_hash for a Torrent.

When a node receives a find node query, it should respond with a key "nodes"
containing information for the target node, or the closest K nodes to the
target.

```shell
$ dhtcli query find_node router.utorrent.com:6881 E2467CBF021192C241367B892230DC1E05C058EE
{
  "t": "0x54efbfbd",
  "y": "r",
  "r": {
    "id": "0xebff36697351ff4aec29cdbaabf2fbe3467cc267",
    "nodes": [
      {
        "id": "0x50ae3167e4af7e2488f71039067cc206f5445a02",
        "address": "24.211.42.62:46817"
      },
      {
        "id": "0xaaece75eb0c8c5b03506e6ff76fb186f531bd525",
        "address": "122.22.198.41:11177"
      },
      {
        "id": "0xd5fcfd487304880469dd06aa52cb82ac2aeaf002",
        "address": "36.230.61.145:16001"
      },
      {
        "id": "0xbea2c90a99ab26d7cd7ca7b25e9d0b729ea5da8d",
        "address": "115.124.173.169:23113"
      },
      {
        "id": "0x8fd573ca5252f31b8fd64307ccac7ed31950f5c2",
        "address": "64.189.196.117:45093"
      },
      {
        "id": "0x3c038cebe963ef32123e1eea127bf5d6925b664e",
        "address": "191.177.186.82:16086"
      },
      {
        "id": "0xed919650393c8090cd3df3d94b78e1ddba10f2cc",
        "address": "95.25.75.150:24296"
      },
      {
        "id": "0x707a0800bd28eeba6fb5a3cfd9bfeeecc2817543",
        "address": "222.99.238.24:29428"
      },
      {
        "id": "0x3113a4fb0af545e489fcf5d6bd117d40f7d4d3bc",
        "address": "2.184.216.174:35462"
      }
    ]
  },
  "v": "0x"
}
```

#### get_peers

Get peers associated with a torrent info_hash from a DHT node.

If the queried node has peers for the info_hash, they are returned in a key
"values" as a list of strings. Values are in IP:Port format. If the queried
node has no peers for the info_hash, a key "nodes" is returned containing the K
closest nodes to the info_hash.  In either case, a "token" key is also included
in the return value. This token value is required for a future announce_peer
query.

```shell
$ dhtcli query get_peers 79.22.72.60:55455 E2467CBF021192C241367B892230DC1E05C0580E
{
  "t": "0x0d4e",
  "y": "r",
  "r": {
    "id": "0xe246e7edc7cf68b7f3abf20ddb042f8182312e0e",
    "nodes": [
      {
        "id": "0xe2462bdd340f0ae1b51a7b1d98f3ad8e6e33e7ac",
        "address": "189.101.214.57:13358"
      },
      {
        "id": "0xe24678d6ae529049f1f1bbe9ebb3a6db3c870ce1",
        "address": "85.66.198.68:5794"
      },
      {
        "id": "0xe2460f49f1f1bbe9ebb3a6db3c870c3e99245e52",
        "address": "95.149.5.53:42061"
      },
      {
        "id": "0xe246555b4be3cf1d716c845218e9fe54feaa9c65",
        "address": "101.165.148.218:18437"
      },
      {
        "id": "0xe246405e58e82ba515172f8d0e2e69d00ea88a63",
        "address": "83.16.219.194:6889"
      }
    ],
    "p": "0x961f",
    "token": "0x29e67ee7",
    "values": [
      "39.8.43.112:25080",
      "41.83.3.125:23227",
      "47.213.57.140:23985",
      "59.46.188.197:19103",
      "68.133.23.82:52481",
      "78.56.91.69:49369",
      "86.210.223.121:31316",
      "90.146.77.100:9225",
      "94.212.49.150:59831",
      "103.106.150.134:43895",
      "107.15.138.17:37120",
      "166.173.62.88:44149",
      "179.186.148.131:9668",
      "197.148.1.17:49291"
    ]
  },
  "v": "0x4c540100"
}
```

#### announce_peer

Announce ourselves to a DHT node as a peer for the torrent with the given
info_hash.

This request requires a token received from the node in a previous get_peers
request. This token is used to prevent malicious hosts from signing up other
hosts for torrents. If --token is not specified, one will be obtained by
issuing a get_peers request to the node. --port specifies the port of the
announced peer. If port is set to 0, or unset, the announce_peer request will
contain the "implied_port" setting. This setting will derive the port value
automatically as described in BEP 5.

```shell
$ dhtcli query announce_peer 79.22.72.60:55455 E2467CBF021192C241367B892230DC1E05C0580E
2019/11/15 19:27:36 --token not specified, issuing get_peers request first to obtain one.
2019/11/15 19:27:37 Got token 0x29e67ee7.
{
  "t": "0x63efbfbd",
  "y": "r",
  "r": {
    "id": "0xe246e7edc7cf68b7f3abf20ddb042f8182312e0e",
    "p": "0x826c"
  },
  "v": "0x4c540100"
}
```

### DHT (experimental)

Issues full requests to the BitTorrent DHT.

#### find_node

Find node is used to find the contact information for a node given its ID.

Response will contain a key "nodes" containing information for the target node
and/or the closest K nodes to the target. Response will be from the closest
node to the target.

```shell
$ dhtcli dht find_node E2467CBF021192C241367B892230DC1E05C0580E
{
  "t": "0xefbfbdefbfbd",
  "y": "r",
  "r": {
    "id": "0xe2467dd3b3f9e37759405aad2d4dce177c9dad26",
    "nodes": [
      {
        "id": "0xe2467a0d2dfd90ff34037eea837733a5d714413a",
        "address": "76.107.99.114:40959"
      },
      {
        "id": "0xe24678d6ae529049f1f1bbe9ebb3a6db3c870ce1",
        "address": "85.66.198.68:5794"
      },
      {
        "id": "0xe24673785c2a15fddd215da5c6f2e3bcea97963b",
        "address": "187.37.135.214:8999"
      },
      {
        "id": "0xe2466ee77a4a205324faf46028cd7a7a7b24d220",
        "address": "71.222.58.88:50321"
      },
      {
        "id": "0xe2466e3fa9d1690b17e6bcaaaff9e0e4ce438ef4",
        "address": "69.59.40.78:6881"
      }
    ],
    "p": "0xa336"
  },
  "v": "0x4c540102"
}
```

Notice how node id's are "close by" to the ID of the target parameter.
