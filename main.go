package main

import (
	"fmt"
	"os"
	"regexp"
	"sync"

	configextractor "bartering/config-extractor"
	datastructures "bartering/data-structures"
	fswatcher "bartering/fs-watcher"
	"bartering/functions"
	peersconnect "bartering/peers-connect"
	storagetesting "bartering/storage-testing"
)

var NodeStorage float64

var port = "8081"

func main() {
	msgCounter := 0

	args := os.Args

	ipRegex := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)

	if len(args) != 2 {
		fmt.Println("Not enough arguments ; use : ./bartering <bootstrap-IP>")
		os.Exit(1)
	} else if !ipRegex.MatchString(args[1]) {
		fmt.Println("Argument invalid : must be an IP address")
		os.Exit(1)
	}

	bootstrapIp := args[1]

	fmt.Println("Extracting configuration")
	config := configextractor.ConfigExtractor("config.yaml")
	configextractor.ConfigPrinter(config)

	storagePool, pendingRequests, fulfilledRequests, peers, bytesAtPeers, bytesForPeers, scores, ratiosAtPeers, ratiosForPeers, storedForPeers := functions.NodeStartup(bootstrapIp)

	fmt.Println("Bytes at peers :", bytesAtPeers)
	fmt.Println("Bytes stored for peers : ", bytesForPeers)
	fmt.Println("Fulfilled requests : ", fulfilledRequests)
	fmt.Println("Storage pool : ", storagePool)
	fmt.Println("Pending requests : ", pendingRequests)
	fmt.Println("Peers : ", peers)
	fmt.Println("Scores : ", scores)
	fmt.Println("Node ratios : ", ratiosForPeers)
	fmt.Println("ratios at peers : ", ratiosAtPeers)
	fmt.Println("stored for peers : ", storedForPeers)
	fmt.Println("")
	fmt.Println("Node started ! Listening on port ", port)

	decreaseBehavior, increaseBehavior := functions.IncreaseDecreaseBehaviors(config)

	var wg sync.WaitGroup
	deletionQueue := []datastructures.StorageRequestTimedAccepted{}

	wg.Add(1)
	go func() {
		// Peer listener: receive and answer messages from other peers.
		defer wg.Done()
		peersconnect.ListenPeersRequestsTCP(port, NodeStorage, bytesAtPeers, scores, ratiosAtPeers, ratiosForPeers, bytesForPeers, &storedForPeers, config.BarteringFactorAcceptableRatio, &deletionQueue, &msgCounter)
	}()

	wg.Add(1)
	go func() {
		// Storage testing: periodically request proof of storage from peers.
		defer wg.Done()
		storagetesting.PeriodicTests(&fulfilledRequests, scores, config.StoragetestingTimerTimeoutSec, port, config.StoragetestingTestingPeriod, decreaseBehavior, increaseBehavior, bytesAtPeers, config.StoragerequestsScoreDecreaseRefusedStoReq)
	}()

	wg.Add(1)
	go func() {
		// Filesystem watcher: replicate any file added to ./data onto the network.
		defer wg.Done()
		fswatcher.FsWatcher("./data", scores, config.DataCopies, port, bytesAtPeers, &fulfilledRequests, config.StoragerequestsScoreDecreaseRefusedStoReq)
	}()

	wg.Wait()
}
