package functions

import (
	"fmt"

	bootstrapconnect "bartering/bootstrap-connect"
	configextractor "bartering/config-extractor"
	datastructures "bartering/data-structures"
)

func NodeStartup(bootstrapIp string) ([]string, []datastructures.StorageRequest, []datastructures.FulfilledRequest, []string, []datastructures.PeerStorageUse, []datastructures.PeerStorageUse, []datastructures.NodeScore, []datastructures.NodeRatio, []datastructures.NodeRatio, []datastructures.FulfilledRequest) {

	/*
		Function called upon a node's startup
		This function will create all needed lists :
			- storagePool : list of CIDs of node's data
			- pendingRequests : list of storage requests that are awaiting to be given to a peer in the network
			- fulfilledRequests : requests accepted by other peers in the network
			- peers : list of peers ids
			- bytesAtPeers : list of bytes stored by other peers
			- scores : list of scores attributed to each peer
			- ratios : storage ratios for each peer

		This function will call the bootstrap to retrieve peers' IP addresses, and store them in the peers list

		Arguments : None
		Returns : storagePool as list of strings, pendingRequests and fulfilledRequests as StorageRequests lists, peers as list of strings
	*/

	fmt.Println("Starting node")

	fmt.Println("Creating storage pool and requests lists")
	storage_pool, pending_requests, _ := createStorageRequestsLists()

	fulfilled_requests := []datastructures.FulfilledRequest{}

	fmt.Println("Creating peers list")
	peers := bootstrapconnect.GetPeersFromBootstrapHTTP(bootstrapIp, "8082")

	fmt.Println("Creating bytes at peers list, scores and ratios")
	bytesAtPeers := initiatePeerStorageUseArray(peers, 0.0)
	bytesForPeers := initiatePeerStorageUseArray(peers, 0.0)
	scores := initiateScores(peers, 10.0)
	ratiosForPeers := initiateRatios(peers, 1.0)
	ratiosAtPeers := initiateRatios(peers, 1.0)

	storedForPeers := []datastructures.FulfilledRequest{}

	return storage_pool, pending_requests, fulfilled_requests, peers, bytesAtPeers, bytesForPeers, scores, ratiosAtPeers, ratiosForPeers, storedForPeers
}

func IncreaseDecreaseBehaviors(config configextractor.Config) ([]datastructures.ScoreVariationScenario, []datastructures.ScoreVariationScenario) {
	DecreasingBehavior := []datastructures.ScoreVariationScenario{{Scenario: "failedTestTimeout", Variation: config.StoragetestingFailedTestTimeoutDecrease}, {Scenario: "failedTestWrongAns", Variation: config.StoragetestingFailedTestWrongAnsDecrease}}

	IncreasingBehavior := []datastructures.ScoreVariationScenario{{Scenario: "passedTest", Variation: config.StoragetestingPassedTestIncrease}}

	return DecreasingBehavior, IncreasingBehavior
}

func initiatePeerStorageUseArray(peers []string, initialStorage float64) []datastructures.PeerStorageUse {

	/*
		Function to initiate an array of PeerStorageUse objects
		Arguments : list of peers, initialStorage value (usually at 0.0)
		Output : array of PeerStorageUse objects
	*/

	bytesAtPeers := []datastructures.PeerStorageUse{}
	for _, peer := range peers {
		bytesAtPeers = append(bytesAtPeers, datastructures.PeerStorageUse{NodeIP: peer, StorageAtNode: initialStorage})
	}
	return bytesAtPeers
}

func initiateScores(peers []string, initialScore float64) []datastructures.NodeScore {

	/*
		Function to iniate an array of bartering.NodeScore objects
		Arguments : list of peers, initial score value
		Output : bartering.NodeScore array
	*/

	scores := []datastructures.NodeScore{}
	for _, peer := range peers {
		scores = append(scores, datastructures.NodeScore{NodeIP: peer, Score: initialScore})
	}
	return scores
}

func initiateRatios(peers []string, initialRatio float64) []datastructures.NodeRatio {

	/*
		Function to initiate array of bartering.NodeRatio objects
		Arguments : list of peers, initiat ratio value
		Output : bartering.NodeRatio array
	*/

	ratios := []datastructures.NodeRatio{}
	for _, peer := range peers {
		ratios = append(ratios, datastructures.NodeRatio{NodeIP: peer, Ratio: initialRatio})
	}
	return ratios
}

func createStorageRequestsLists() ([]string, []datastructures.StorageRequest, []datastructures.StorageRequest) {

	/*
		Function to create all needed data structures
		Argument : None
		Returns : storage_pool as string list, pending and fulfilled requests lists as StorageRequest lists
	*/

	storage_pool := []string{}

	pending_requests := []datastructures.StorageRequest{}

	fulfilled_requests := []datastructures.StorageRequest{}

	return storage_pool, pending_requests, fulfilled_requests

}
