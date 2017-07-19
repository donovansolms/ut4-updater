package ut4updater_test

import (
	"testing"

	"github.com/donovansolms/ut4-updater"
)

var updater *ut4updater.UT4Updater

func TestMain(m *testing.M) {
	var err error
	updater, err = ut4updater.New(
		"./test-resources/installs",
		2,
		"latest",
		true,
		"http://localhost/ut4updater/versionmap.json")
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestGetVersionList(t *testing.T) {
	versions, err := updater.GetVersionList()
	if err != nil {
		t.Error(err.Error())
	}
	if len(versions) == 0 {
		t.Error("Versions list is empty")
	}
}

func TestGetLatestVersion(t *testing.T) {
	latestVersion, err := updater.GetLatestVersion()
	if err != nil {
		t.Error(err.Error())
	}
	if latestVersion.Version != "003" {
		t.Errorf("Returned version '%s'. Expected '%s'",
			latestVersion.Version,
			"003")
	}
}
