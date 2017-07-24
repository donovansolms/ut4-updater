package ut4updater

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var updater *UT4Updater

func TestMain(m *testing.M) {
	var err error
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() == "/update/ut4-versionmap" {
			versionMap, err := ioutil.ReadFile("./test-resources/installs/versionmap.json")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(versionMap)
		} else if r.URL.EscapedPath() == "/update/ut4-check" {
			// TODO: Implement
			w.WriteHeader(http.StatusOK)
		} else if r.URL.EscapedPath() == "/update/ut4-hash/latest" {
			w.Write([]byte("{\"Unreal.pak\": \"1234567890oiuytrewq\"}"))
		} else if r.URL.EscapedPath() == "/update/ut4-update/deb3e700df1e6b29df98c26cc388417072b0bb5eeda3de7d035e186c315f161c" {
			w.Write([]byte("{\"update_url\": \"https://www.google.com\"}"))
		}
		//fmt.Println(r.URL.EscapedPath())
	}))
	defer testServer.Close()

	updater, err = New(
		"./test-resources/installs",
		2,
		"latest",
		true,
		testServer.URL)
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

func TestUpdateVersionMap(t *testing.T) {
	// Test file the remote download and local cache
	previousPath := updater.installPath
	previousURL := updater.updateURL
	updater.installPath = "/tmp"
	updater.updateURL = "httx://localhost"
	err := updater.updateVersionMap()
	if err == nil {
		t.Error("Invalid version URL and local file must fail")
	}
	updater.installPath = previousPath
	updater.updateURL = previousURL
}

func TestGetFilelist(t *testing.T) {
	list, err := updater.getFilelist("./test-resources")
	if err != nil {
		t.Error(err.Error())
	}
	if len(list) == 0 {
		t.Errorf("getFilelist must not return an empty list")
	}
}

func TestGenerateHashes(t *testing.T) {
	list, err := updater.getFilelist("./test-resources/installs")
	if err != nil {
		t.Error(err.Error())
	}
	feedbackChan := make(chan HashProgressEvent)
	go updater.GenerateHashes(list, 1, feedbackChan)
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

	hashes, err := updater.getRemoteVersionHashes("latest")
	if err != nil {
		t.Error(err.Error())
	}
	if len(hashes) == 0 {
		t.Error("Remote version hashes must not be empty")
	}
}

func TestCalculateDelta(t *testing.T) {

	current := make(map[string]string)
	next := make(map[string]string)

	current["a"] = "unmodified"
	current["b"] = "modify"
	current["c"] = "c"
	next["a"] = "unmodified"
	next["b"] = "modified"
	//next["c"] = "removed"
	next["d"] = "added"

	deltaOperations := updater.calculateHashDeltaOperations(current, next)
	if len(deltaOperations) == 0 {
		t.Error("Deltas should not be empty")
	}
}

func TestGenerateDeltaHash(t *testing.T) {

	deltaOperations := make(map[string]string)
	deltaOperations["d"] = "added"
	deltaOperations["b"] = "modified"
	deltaOperations["c"] = "removed"

	hash := updater.generateDeltaHash(deltaOperations)
	if hash == "" {
		t.Error("Hash may not be empty")
	}
	if hash != "deb3e700df1e6b29df98c26cc388417072b0bb5eeda3de7d035e186c315f161c" {
		t.Error("Hash doesn't match input data")
	}
}

func TestGetUpdatePackageURL(t *testing.T) {

	versionHash := "deb3e700df1e6b29df98c26cc388417072b0bb5eeda3de7d035e186c315f161c"
	updateURL, err := updater.getUpdatePackageURL(versionHash)
	if err != nil {
		t.Error(err.Error())
	}
	if updateURL == "" {
		t.Error("UpdateURL must not be blank")
	}
}
