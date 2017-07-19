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
		"http://localhost/ut4updater/versionmap.json")
	if err != nil {
		panic(err)
	}
	// Test Fail the remote download
	err = testUpdater.updateVersionMap("httx://localhost/versionmaps.json")
	if err == nil {
		t.Error("Invalid version URL must fail")
	}
	// Test file the remote download and local cache
	testUpdater.installPath = "/tmp"
	err = testUpdater.updateVersionMap("httx://localhost/versionmaps.json")
	if err == nil {
		t.Error("Invalid version URL and local file must fail")
	}
}
