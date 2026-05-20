# Developer guide

This guide covers setting up a development environment, building, running, and
testing pr-bartering, and running several nodes locally.

## Environment

You need:

- Go 1.26 or newer (`go version`). The module pins the toolchain in `go.mod`, so
  with `GOTOOLCHAIN=auto` (the default) the right version is fetched on demand.
- IPFS Kubo on your `PATH` (`ipfs version`). Start it with `ipfs daemon`.
- `make`.

Clone and verify the toolchain:

```
git clone git@github.com:Kitio-Tek/pr-bartering.git
cd pr-bartering
make build
```

## Make targets

```
make build      # compile bin/bartering and bin/bootstrap
make test       # unit tests
make test-race  # unit tests with the race detector
make cover      # tests with a coverage profile, prints total coverage
make fmt        # gofmt the tree
make vet        # go vet
make lint       # golangci-lint, pinned to the version in the Makefile
make sec        # gosec static security scan
make vuln       # govulncheck against deps and the standard library
make check      # fmt, vet, lint, test, sec, vuln
make clean      # remove bin/ and coverage output
```

`make lint` uses the exact golangci-lint version that CI uses, so a local run
reproduces the pipeline.

## Build and run

The node and the bootstrap service are two separate binaries.

```
make build
```

Run an IPFS daemon, then the bootstrap, then a node, each in its own terminal:

```
ipfs daemon
```

```
cd bootstrap-node && ../bin/bootstrap 127.0.0.1   # HTTP on :8082
```

```
./bin/bartering 127.0.0.1                          # TCP on :8081, reads ./config.yaml
```

Add a file to `data/` and the node will replicate it:

```
echo "hello" > data/hello.txt
```

## Running several nodes locally

Each node needs its own working directory (for its own `config.yaml` and `data/`)
and its own listen port. The bootstrap returns every address in
`bootstrap-node/ips.txt` except the caller's, so list the node addresses there.

A minimal two-node setup using loopback aliases:

1. Add loopback aliases (Linux):

   ```
   sudo ip addr add 127.0.0.2/8 dev lo
   ```

2. Put both addresses in `bootstrap-node/ips.txt`, one per line:

   ```
   127.0.0.1
   127.0.0.2
   ```

3. Create a working directory per node, each with a `config.yaml` (copy the one
   at the repository root and change `Port`) and an empty `data/` directory.

4. Start the bootstrap, then start each node from its own directory, passing the
   bootstrap IP.

The node prints the peer list, scores, and ratios it learned at startup, which is
the quickest way to confirm discovery worked.

## Debugging and profiling

The daemon logs to standard output; run a single node in the foreground to follow
the message flow. For CPU or memory profiles, build with the standard `net/http/pprof`
approach in a throwaway branch, or run targeted benchmarks with
`go test -run NONE -bench . -cpuprofile cpu.out ./<package>`.

## Coding standards

- Run `make check` before opening a pull request.
- Keep `gofmt` clean; CI fails on unformatted files.
- New error returns must be handled, not discarded, except where an explicit
  `_ =` documents an intentional ignore.
- Subprocess calls and file reads that gosec flags must either be made safe or
  annotated with a `#nosec` comment that justifies why the finding does not apply.
