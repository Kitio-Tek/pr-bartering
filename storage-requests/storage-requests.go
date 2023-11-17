package storagerequests

import (
	api_ipfs "bartering/api-ipfs"
	bartering "bartering/bartering-api"
	"bartering/utils"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type StorageRequest struct {
	/*
		Data structure to represent storage requests ; consist of the CID of a file and its size
	*/

	CID      string
	fileSize float64
}

type StorageRequestTimed struct {
	CID             string
	FileSize        float64
	DurationMinutes int
}

type FilesAtPeers struct {
	Peer  string
	Files []string
}

func makeStorageRequest() {

}

func HandleStorageRequest(bufferString string) {
	/*
		Function to handle a storage message type message
		Arguments : buffer received through a tcp connection, as a string
	*/

	fmt.Println("Received storage request")
	CID := bufferString[5:51]
	fmt.Println("CID : ", CID)
	fileSize := bufferString[51:]
	fileSize = strings.Split(fileSize, "\n")[0]
	fmt.Println("File Size : ", fileSize)

	fileSizeFloat, err := strconv.ParseFloat(fileSize, 64)
	utils.ErrorHandler(err)

	request := StorageRequest{fileSize: fileSizeFloat, CID: CID}

	fmt.Println("Storage request : ", request, " ; checking validity ...")

	if CheckRqValidity(request) {
		fmt.Println("Request ", request, " valid, storing ! ")
		fmt.Println("Pinning to IPFS ...")
		api_ipfs.PinToIPFS(CID)
		fmt.Println("File pinned to IPFS!")
	} else {
		fmt.Println("Request ", request, " not valid, not storing ! ")
	}

}

func checkRatioValidity(peer string, ratios []bartering.NodeRatio, bytesAtPeers []bartering.PeerStorageUse, storedFor []bartering.PeerStorageUse) {
	// ratio, err := bartering.FindNodeRatio(ratios, peer)
	// utils.ErrorHandler(err)

}

func watchPendingRequests() {

}

func ElectStorageNodes(peerScores []bartering.NodeScore, numberOfNodes int) ([]string, error) {
	/*
		Function to elect nodes to whom self will send storage requests
		Arguments : IP of peer as string, number of nodes as int
		Returns : list of strings containing IPs of nodes to contact
	*/

	// TODO change so we also elect low score nodes to give the opportunity to raise the score

	if numberOfNodes > len(peerScores) {
		return []string{}, errors.New("asking for more peers than we know")
	}

	electedNodesScores := []bartering.NodeScore{}
	for _, peerScore := range peerScores {

		if len(electedNodesScores) < numberOfNodes {

			electedNodesScores = append(electedNodesScores, peerScore)

		} else {
			for index, currentlyElectedNode := range electedNodesScores {

				if currentlyElectedNode.Score < peerScore.Score {
					electedNodesScores[index] = peerScore
					break
				}
			}
		}

	}

	electedNodes := []string{}

	for _, electedNodeScore := range electedNodesScores {
		electedNodes = append(electedNodes, electedNodeScore.NodeIP)
	}

	return electedNodes, nil
}

func CheckRqValidity(storageRequest StorageRequest) bool {

	return true
}

func CheckCIDValidity(storageRequest StorageRequest) bool {
	/*
		Check if : CID is valid and exists
		problem : ipfs cat with wrong CID goes through ipfs search which can be very long
		no efficient way to check if CID exists
		however if we work under the hypothesis that peers in our network have preestablished ipfs links
	*/
	return true
}

func CheckFileSizeValidity(storageRequest StorageRequest) bool {
	/*
		Check if fileSize announced in storage request is declared honestly
	*/
	return true
}

func CheckEnoughSpace(storageRequest StorageRequest, currentStorageSpace float64, NodeTotalStorageSpace float64) bool {
	/*
		Check if self has enough space to store the file
		Arguments : storage request of type StorageRequest, current storage space used as float64
		Returns : boolean
	*/

	return storageRequest.fileSize+float64(currentStorageSpace) < NodeTotalStorageSpace

}
