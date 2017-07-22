package ut4updater_test

import (
	"os"
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
		"http://update.donovansolms.local")
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
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

func TestGetOSDistribution(t *testing.T) {
	osDistribution := updater.GetOSDistribution()
	if osDistribution.Distribution == "" {
		t.Error("Distribution must contain something")
	}
	if osDistribution.DistributionID == "" {
		t.Error("DistributionID must contain something")
	}
	if osDistribution.KernelVersion == "" {
		t.Error("KernelVersion must contain something")
	}
}

func TestUpdateCheck(t *testing.T) {
	shouldUpdate, err := updater.CheckForUpdate()
	if err != nil {
		t.Error(err.Error())
	}
	if shouldUpdate == false {
		t.Error("UpdateCheck should return true")
	}
}
