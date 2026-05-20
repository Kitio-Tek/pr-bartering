package functions

import (
	datastructures "bartering/data-structures"
	"bartering/utils"
	"os"
	"path/filepath"
	"testing"
)

func TestGetFileSize(t *testing.T) {
	// A 13-byte file should report as 13/1024 KB.
	path := filepath.Join(t.TempDir(), "fixture.txt")
	if err := os.WriteFile(path, []byte("hello, world!"), 0o600); err != nil {
		t.Fatalf("could not write fixture: %v", err)
	}

	result := utils.GetFileSize(path)
	expected := 13.0 / 1024.0
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestInitiateBytesAtPeers(t *testing.T) {
	peers := []string{"peer1", "peer2"}
	storageAtPeer1 := datastructures.PeerStorageUse{NodeIP: "peer1", StorageAtNode: 0.0}
	storageAtPeer2 := datastructures.PeerStorageUse{NodeIP: "peer2", StorageAtNode: 0.0}
	result := initiatePeerStorageUseArray(peers, 0.0)
	if result[0] != storageAtPeer1 || result[1] != storageAtPeer2 {
		t.Errorf("BytesAtPeers not initiated correctly")
	}
}

func TestInitiateScores(t *testing.T) {
	peers := []string{"peer1", "peer2"}
	peerScore1 := datastructures.NodeScore{NodeIP: "peer1", Score: 10.0}
	peerScore2 := datastructures.NodeScore{NodeIP: "peer2", Score: 10.0}

	result := initiateScores(peers, 10.0)

	if result[0] != peerScore1 || result[1] != peerScore2 {
		t.Errorf("Scores not initiated correctly")
	}
}
