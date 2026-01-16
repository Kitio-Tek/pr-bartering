package storagerequests

import (
	datastructures "bartering/data-structures"
	"fmt"
	"testing"
	"time"
)

func TestAuxInsertInSortedList(t *testing.T) {

	queue := []datastructures.StorageRequestTimedAccepted{}
	storageRequestFirst := datastructures.StorageRequestTimedAccepted{CID: "blablablafirst", Deadline: time.Now()}

	i := 0
	for i < 3 {
		storageRequest := datastructures.StorageRequestTimedAccepted{CID: "blablabla" + fmt.Sprint(i), Deadline: time.Now()}
		queue = append(queue, storageRequest)
		i += 1
	}
	time.Sleep(5 * time.Second)
	storageRequestLast := datastructures.StorageRequestTimedAccepted{CID: "blablablalast", Deadline: time.Now()}

	newQueue := AuxInsertInSortedList(storageRequestFirst, queue)
	newNewQueue := AuxInsertInSortedList(storageRequestLast, newQueue)

	if newNewQueue[0] != storageRequestFirst || newNewQueue[len(newNewQueue)-1] != storageRequestLast {
		t.Errorf("timed storage requests not inserted correctly into deletion queue ")
	}

}

func TestAppendStorageRequestToDeletionQueue(t *testing.T) {
	queue := []datastructures.StorageRequestTimedAccepted{}
	storageRequestFirst := datastructures.StorageRequestTimedAccepted{CID: "blablablafirst", Deadline: time.Now()}
	i := 0
	for i < 3 {
		storageRequest := datastructures.StorageRequestTimedAccepted{CID: "blablabla" + fmt.Sprint(i), Deadline: time.Now()}
		queue = append(queue, storageRequest)
		i += 1
	}
	fmt.Println(queue)
	time.Sleep(5 * time.Second)
	storageRequestLast := datastructures.StorageRequestTimedAccepted{CID: "blablablalast", Deadline: time.Now()}

	AppendStorageRequestToDeletionQueue(storageRequestFirst, &queue)
	AppendStorageRequestToDeletionQueue(storageRequestLast, &queue)

	fmt.Println(queue)

	if queue[0] != storageRequestFirst || queue[len(queue)-1] != storageRequestLast {
		t.Errorf("timed storage requests not inserted correctly into deletion queue")
	}

}

func TestGarbageCollectionStrategy(t *testing.T) {

	storageDeletionQueue := []datastructures.StorageRequestTimedAccepted{}
	i := 0
	for i < 3 {
		storageRequest := datastructures.StorageRequestTimedAccepted{CID: "blablabla" + fmt.Sprint(i), Deadline: time.Now()}
		storageDeletionQueue = append(storageDeletionQueue, storageRequest)
		i += 1
	}

	for i >= 0 {
		storageDeletionQueue = GarbageCollectionStrategy(storageDeletionQueue)
		i -= 1
	}

	if len(storageDeletionQueue) != 0 {
		t.Errorf("garbage collection strategy not behaving properly")
	}

}

func TestElectStorageNodesLowAndHigh(t *testing.T) {
	scores := []datastructures.NodeScore{}
	i := 0
	score := 3.0
	for i < 30 {
		scores = append(scores, datastructures.NodeScore{
			NodeIP: fmt.Sprintf("127.0.0.%d", i+1), 
			Score:  score,
		})
		score += 0.5
		i++
	}

	elected10 := ElectStorageNodesLowAndHigh(scores, 10)
	fmt.Println("Elected 10 nodes (low and high):", elected10)
	if len(elected10) != 10 {
		t.Errorf("Expected 10 nodes, got %d", len(elected10))
	}

	elected20 := ElectStorageNodesLowAndHigh(scores, 20)
	fmt.Println("Elected 20 nodes (low and high):", elected20)
	if len(elected20) != 20 {
		t.Errorf("Expected 20 nodes, got %d", len(elected20))
	}
}
// func TestComputeDeadlineFromTimedStorageRequest(t *testing.T) {
// 	storageRequest := datastructures.StorageRequestTimed{CID: "whatever", DurationMinutes: 3}
// 	currentTime := time.Now()
// 	deadline := ComputeDeadlineFromTimedStorageRequest(storageRequest)
// 	fmt.Println(deadline.Sub(currentTime))

// 	if deadline.Sub(currentTime)-time.Duration(storageRequest.DurationMinutes) > time.Duration(0.00005) {
// 		t.Errorf("garbage collection strategy not behaving properly")
// 	}
// }

func TestElectStorageNodes(t *testing.T) {
	scores := []datastructures.NodeScore{}
	i := 0
	score := 3.0
	for i < 30 {
		// Create unique IP addresses for each node
		scores = append(scores, datastructures.NodeScore{
			NodeIP: fmt.Sprintf("127.0.0.%d", i+1), 
			Score:  score,
		})
		score += 0.5
		i++
	}

	
	elected, err := ElectStorageNodes(scores, 10, make(map[string]bool))
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(elected) != 10 {
		t.Errorf("Expected 10 nodes, got %d", len(elected))
	}
	fmt.Println("Elected 10 nodes:", elected)

	
	usedPeers := make(map[string]bool)
	usedPeers["127.0.0.1"] = true
	usedPeers["127.0.0.2"] = true
	
	elected2, err := ElectStorageNodes(scores, 10, usedPeers)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(elected2) != 10 {
		t.Errorf("Expected 10 nodes, got %d", len(elected2))
	}
	for _, peer := range elected2 {
		if usedPeers[peer] {
			t.Errorf("Used peer %s was elected", peer)
		}
	}
	fmt.Println("Elected 10 nodes (excluding used):", elected2)

	_, err = ElectStorageNodes(scores, 35, make(map[string]bool))
	if err == nil {
		t.Errorf("Expected error when asking for more nodes than available")
	}


	manyUsedPeers := make(map[string]bool)
	for i := 0; i < 25; i++ {
		manyUsedPeers[fmt.Sprintf("127.0.0.%d", i+1)] = true
	}
	_, err = ElectStorageNodes(scores, 10, manyUsedPeers)
	if err == nil {
		t.Errorf("Expected error when not enough peers available after exclusions")
	}

	allElected, err := ElectStorageNodes(scores, 30, make(map[string]bool))
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(allElected) != 30 {
		t.Errorf("Expected 30 nodes, got %d", len(allElected))
	}
	fmt.Println("All 30 nodes elected successfully")
}

