package peersconnect

import (
	"fmt"
	"net"
	"time"

	"bartering/bartering-api"
	datastructures "bartering/data-structures"
	storagerequests "bartering/storage-requests"
	storagetesting "bartering/storage-testing"
	"bartering/utils"
)

func ListenPeersRequestsTCP(port string, nodeStorage float64, bytesAtPeers []datastructures.PeerStorageUse, scores []datastructures.NodeScore, ratiosAtPeers []datastructures.NodeRatio, ratiosForPeers []datastructures.NodeRatio, bytesForPeers []datastructures.PeerStorageUse, storedForPeers *[]datastructures.FulfilledRequest, factorAcceptableRatio float64, deletienQueue *[]datastructures.StorageRequestTimedAccepted, msgCounter *int) {

	/*
		TCP server to receive messages from peers
	*/
	maxConnections := 4 // Maximum number of concurrent connections
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		utils.ErrorHandler(err)
		return
	}
	defer func() { _ = listener.Close() }()

	// Create a channel to limit the number of concurrent connections
	connLimiter := make(chan struct{}, maxConnections)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		connLimiter <- struct{}{} // Block if the limit is reached
		go func(conn net.Conn) {
			defer func() {
				<-connLimiter // Release the slot
				_ = conn.Close()
			}()

			// Set TCP keep-alive
			if tcpConn, ok := conn.(*net.TCPConn); ok {
				_ = tcpConn.SetKeepAlive(true)
				_ = tcpConn.SetKeepAlivePeriod(3 * time.Minute) // Keep-alive period
			}

			handleConnection(conn, nodeStorage, bytesAtPeers, scores, ratiosAtPeers, bytesForPeers, storedForPeers, factorAcceptableRatio, deletienQueue, msgCounter)
		}(conn)
	}
}

func handleConnection(conn net.Conn, nodeStorage float64, bytesAtPeers []datastructures.PeerStorageUse, scores []datastructures.NodeScore, ratios []datastructures.NodeRatio, bytesForPeers []datastructures.PeerStorageUse, storedForPeers *[]datastructures.FulfilledRequest, factorAcceptableRatio float64, deletionQueue *[]datastructures.StorageRequestTimedAccepted, msgCounter *int) {

	/*
		Connection handler for TCP connections received through the TCP server
		Arguments : a connection as net.Conn
	*/

	defer func() { _ = conn.Close() }()

	buffer := make([]byte, 63)

	if _, err := conn.Read(buffer); err != nil {
		return
	}

	MessageDiscriminator(buffer, conn, nodeStorage, bytesAtPeers, scores, ratios, bytesForPeers, storedForPeers, factorAcceptableRatio, deletionQueue, msgCounter)

}

func MessageDiscriminator(buffer []byte, conn net.Conn, nodeStorage float64, bytesAtPeers []datastructures.PeerStorageUse, scores []datastructures.NodeScore, ratios []datastructures.NodeRatio, bytesForPeers []datastructures.PeerStorageUse, storedForPeers *[]datastructures.FulfilledRequest, factorAcceptableRatio float64, deletionQueue *[]datastructures.StorageRequestTimedAccepted, msgCounter *int) {

	/*
		Function used to discriminate different types of messages and call the necessary functions for each type of messages
		Arguments : a slide of bytes []byte
	*/

	bufferString := string(buffer)
	messageType := bufferString[:5]

	switch messageType {
	case "StoRq":
		storagerequests.HandleStorageRequest(bufferString, conn, bytesForPeers, storedForPeers)

	case "BarRq":
		remoteAddr := conn.RemoteAddr()
		ip, _, err := net.SplitHostPort(remoteAddr.String())
		utils.ErrorHandler(err)
		fmt.Println("Received bartering request from peer", ip)
		bartering.RespondToBarterMsg(bufferString, ip, nodeStorage, bytesAtPeers, scores, conn, ratios, factorAcceptableRatio, msgCounter)

	case "TesRq":
		fmt.Println("Received test request")
		CID := bufferString[5 : len(bufferString)-1]
		fmt.Println(CID)
		storagetesting.HandleTest(CID, conn)

	default:
		fmt.Println("Unrecognized message : ", bufferString)
	}
}
