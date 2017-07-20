package ut4updater

import (
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
	testUpdater.installPath = "/tmp"
	err = testUpdater.updateVersionMap("httx://localhost/versionmaps.json")
	if err == nil {
		t.Error("Invalid version URL and local file must fail")
	}
	testUpdater.installPath = previousPath
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
