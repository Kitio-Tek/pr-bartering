# pr-bartering

[![CI](https://github.com/Kitio-Tek/pr-bartering/actions/workflows/ci.yml/badge.svg)](https://github.com/Kitio-Tek/pr-bartering/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.26-00ADD8.svg)](go.mod)

pr-bartering is an experimental overlay on top of [IPFS](https://ipfs.tech) that
lets a set of nodes replicate each other's data by bartering storage. Instead of
paying for pinning with a cryptocurrency, as services such as Filecoin do, each
node trades space: it stores data for its peers in exchange for them storing its
own, and it keeps a per-peer score and storage ratio that govern how much it is
willing to give.

## Contents

- [How it works](#how-it-works)
- [Status](#status)
- [Prerequisites](#prerequisites)
- [Quick start](#quick-start)
- [Configuration](#configuration)
- [Architecture](#architecture)
- [Development](#development)
- [Testing](#testing)
- [Security](#security)
- [Contributing](#contributing)
- [License](#license)

## How it works

A node watches a local `data/` directory. When a file appears, the node adds it
to its IPFS daemon to obtain a content identifier (CID), then asks peers to pin
that CID until the configured number of copies (`DataCopies`) is reached. Peers
decide whether to accept based on the requesting node's score and the storage
ratio currently agreed between them.

Nodes periodically test the peers holding their data by asking for a proof
(today, a SHA-256 hash of the stored content). A peer that answers correctly
gains score; one that times out or answers incorrectly loses score, and the data
is re-replicated elsewhere.

Peers discover each other through a small bootstrap service that returns the list
of known node addresses over HTTP.

## Status

| Capability | State |
| --- | --- |
| Peer discovery via the bootstrap service | Working |
| Replicating files added to `data/` to `DataCopies` peers | Working |
| Responding to storage, barter, and test requests from peers | Working |
| Periodic proof-of-storage testing and score adjustment | Working |
| Responding to a barter (ratio) negotiation started by a peer | Working |
| A node proactively initiating barter rounds | Library present, not yet driven by the daemon loop |
| Node failure simulation (Weibull / lognormal session models) | Library present, used for experiments |

## Prerequisites

- [Go](https://go.dev/dl/) 1.26 or newer.
- A running [IPFS Kubo](https://docs.ipfs.tech/install/command-line/) daemon. The
  protocol shells out to the `ipfs` binary, which must be on your `PATH`.

Verify both:

```
go version
ipfs version
```

## Quick start

The protocol needs three things running: an IPFS daemon, the bootstrap service,
and at least one node. Use a separate terminal for each.

1. Start IPFS:

   ```
   ipfs daemon
   ```

2. Build the binaries and start the bootstrap service. The bootstrap reads the
   peer addresses it hands out from `bootstrap-node/ips.txt` and listens on port
   8082:

   ```
   make build
   cd bootstrap-node && ../bin/bootstrap 127.0.0.1
   ```

3. In another terminal, start a node, passing the bootstrap IP. The node listens
   on port 8081 and reads its settings from `config.yaml` in the working
   directory:

   ```
   ./bin/bartering 127.0.0.1
   ```

4. Drop a file into `data/`. The node detects it, adds it to IPFS, and asks peers
   to store it:

   ```
   echo "hello bartering" > data/hello.txt
   ```

For running several nodes on one machine, see [DEVELOPER.md](DEVELOPER.md).


## Architecture

```
                 +-------------------+
                 |  bootstrap-node   |  HTTP :8082, serves the peer list
                 +-------------------+
                          ^
                          | peer list
                          |
   data/ ---> +-----------------------+  TCP :8081   +------------------+
   (watched)  |        node           | <----------> |   other nodes    |
              |  fs-watcher           |   StoRq      |                  |
              |  storage-requests     |   BarRq      +------------------+
              |  bartering-api        |   TesRq
              |  storage-testing      |
              +-----------------------+
                          |
                          | ipfs add / pin / cat
                          v
                   +--------------+
                   | IPFS (Kubo)  |
                   +--------------+
```

Peer messages are framed with a five-character type prefix: `StoRq` (storage
request), `BarRq` (barter request), and `TesRq` (test request). 
More detail is in [docs/architecture.md](docs/architecture.md) and the design
decisions are recorded in [docs/adr](docs/adr).

## Development

See [DEVELOPER.md](DEVELOPER.md) for the full workflow. The common targets:

```
make build   # compile the node and bootstrap binaries into bin/
make test    # run the unit tests
make lint    # run golangci-lint (pinned version)
make check   # format, vet, lint, test, and run the security scanners
```

## Testing

```
make test        # unit tests
make test-race   # unit tests with the race detector
make cover       # tests with a coverage profile
```

## Security

To report a vulnerability, follow the process in [SECURITY.md](SECURITY.md).
Please do not open a public issue for security problems.

## Contributing

Contributions are welcome. Read [CONTRIBUTING.md](CONTRIBUTING.md) for the commit
style, the pull-request process, and the testing expectations, and
[CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) for the standards we hold each other to.

## License

pr-bartering is licensed under the Apache License 2.0. See [LICENSE](LICENSE).
