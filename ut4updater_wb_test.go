package ut4updater

import (
	"fmt"
	"testing"
)

func TestUpdateVersionMap(t *testing.T) {
	testUpdater, err := New(
		"./test-resources/installs",
		2,
		"latest",
		true,
		"http://update.donovansolms.local")
	if err != nil {
		panic(err)
	}
	// Test file the remote download and local cache
	previousPath := testUpdater.installPath
	previousURL := testUpdater.updateURL
	testUpdater.installPath = "/tmp"
	testUpdater.updateURL = "httx://localhost"
	err = testUpdater.updateVersionMap()
	if err == nil {
		t.Error("Invalid version URL and local file must fail")
	}
	testUpdater.installPath = previousPath
	testUpdater.updateURL = previousURL
}

func TestGetFilelist(t *testing.T) {
	testUpdater, err := New(
		"./test-resources/installs",
		2,
		"latest",
		true,
		"http://update.donovansolms.local")
	if err != nil {
		panic(err)
	}

	list, err := testUpdater.getFilelist("./test-resources")
	if err != nil {
		t.Error(err.Error())
	}
	if len(list) == 0 {
		t.Errorf("getFilelist must not return an empty list")
	}
}

func TestGenerateHashes(t *testing.T) {
	testUpdater, err := New(
		"./test-resources/installs",
		2,
		"latest",
		true,
		"http://update.donovansolms.local")
	if err != nil {
		panic(err)
	}

	list, err := testUpdater.getFilelist("./test-resources")
	if err != nil {
		t.Error(err.Error())
	}
	feedbackChan := make(chan HashProgressEvent)
	go testUpdater.GenerateHashes(list, 2, feedbackChan)
	completed := 0
	for feedback := range feedbackChan {
		if feedback.Completed {
			completed++
		}
	}
	if len(list) != completed {
		t.Error("Not all hashes were generated for the given list")
	}

}

func TestRemoteVersionHashes(t *testing.T) {
	testUpdater, err := New(
		"./test-resources/installs",
		2,
		"latest",
		true,
		"http://update.donovansolms.local")
	if err != nil {
		panic(err)
	}

	hashes, err := testUpdater.getRemoteVersionHashes("latest")
	if err != nil {
		t.Error(err.Error())
	}
	if len(hashes) == 0 {
		t.Error("Remote version hashes must not be empty")
	}
}

func TestCalculateDelta(t *testing.T) {
	testUpdater, err := New(
		"./test-resources/installs",
		2,
		"latest",
		true,
		"http://update.donovansolms.local")
	if err != nil {
		panic(err)
	}
	current := make(map[string]string)
	next := make(map[string]string)

	current["a"] = "unmodified"
	current["b"] = "modify"
	current["c"] = "c"
	next["a"] = "unmodified"
	next["b"] = "modified"
	//next["c"] = "removed"
	next["d"] = "added"

	deltaOperations := testUpdater.calculateHashDeltaOperations(current, next)
	if len(deltaOperations) == 0 {
		t.Error("Deltas should not be empty")
	}
	fmt.Println(deltaOperations)
}

func TestGenerateDeltaHash(t *testing.T) {
	testUpdater, err := New(
		"./test-resources/installs",
		2,
		"latest",
		true,
		"http://update.donovansolms.local")
	if err != nil {
		panic(err)
	}

	deltaOperations := make(map[string]string)
	deltaOperations["d"] = "added"
	deltaOperations["b"] = "modified"
	deltaOperations["c"] = "removed"

	hash := testUpdater.generateDeltaHash(deltaOperations)
	if hash == "" {
		t.Error("Hash may not be empty")
	}
	if hash != "deb3e700df1e6b29df98c26cc388417072b0bb5eeda3de7d035e186c315f161c" {
		t.Error("Hash doesn't match input data")
	}
}
