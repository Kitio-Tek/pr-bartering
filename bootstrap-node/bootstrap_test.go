package main

import (
	"fmt"
	"testing"
)

func TestBuildPeersIPlist(t *testing.T) {
	// Test with the actual ips.txt file
	ownAddress := "172.12.34.1"
	peers, err := BuildPeersIPlist("ips.txt", ownAddress)
	if err != nil {
		t.Errorf("BuildPeersIPlist failed: %v", err)
		return
	}

	expectedCount := 6
	if len(peers) != expectedCount {
		t.Errorf("Expected %d peers, got %d", expectedCount, len(peers))
	}

	expectedPeers := []string{
		"172.12.34.2",
		"172.12.34.3",
		"172.12.34.4",
		"172.12.34.5",
		"172.12.34.6",
		"172.12.34.7",
	}

	for i, peer := range peers {
		if peer != expectedPeers[i] {
			t.Errorf("Expected peer[%d] to be %s, got %s", i, expectedPeers[i], peer)
		}
	}

	fmt.Println("Peers loaded:", peers)
}
