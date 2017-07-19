package ut4updater

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/sethgrid/pester"
)

const (
	runVersionLatest = "latest"
)

// UT4Updater is the main executor for the updater
type UT4Updater struct {
	installPath  string
	keepVersions uint
	runVersion   string
	sendStats    bool
	updateURL    string
	versionMaps  VersionMaps
	clientID     string
}

// New creates aand initializes a new instance of UT4Updater
func New(installPath string,
	keepVersions uint,
	runVersion string,
	sendStats bool,
	updateURL string) (*UT4Updater, error) {
	updater := &UT4Updater{
		installPath:  installPath,
		keepVersions: keepVersions,
		runVersion:   runVersion,
		sendStats:    sendStats,
		updateURL:    updateURL,
	}
	fullPath, err := filepath.Abs(updater.installPath)
	if err != nil {
		return updater, err
	}
	updater.installPath = fullPath

	versionMapURL := fmt.Sprintf("%s/%s/%s",
		updateURL,
		"update",
		"ut4-versionmap")
	fmt.Println(versionMapURL)
	err = updater.updateVersionMap(versionMapURL)
	if err != nil {
		return updater,
			fmt.Errorf("Unable to update version map '%s': %s",
				versionMapURL,
				err.Error())
	}

	// On the first run we generate a UUID for this client
	clientIDPath := filepath.Join(updater.installPath, ".clientid")
	clientUUID, err := ioutil.ReadFile(clientIDPath)
	if err != nil {
		clientUUID, err := uuid.NewRandom()
		if err != nil {
			return updater, err
		}
		err = ioutil.WriteFile(clientIDPath, []byte(clientUUID.String()), 0644)
		if err != nil {
			return updater, err
		}
	}
	updater.clientID = string(clientUUID)

	return updater, nil
}

// updateVersionMap retrieves the version map from the update server
// and saves a copy locally
func (updater *UT4Updater) updateVersionMap(versionMapURL string) error {

	var mapReader io.ReadCloser
	response, err := http.Get(versionMapURL)
	if err != nil {
		// We were unable to fetch the version map from the remote server
		// now we can check if a local copy exists
		// Declaring localErr to avoid shadowing mapReader
		var localErr error
		mapReader, localErr = os.Open(filepath.Join(
			updater.installPath,
			"versionmap.json"))
		if localErr != nil {
			return fmt.Errorf("Remote returned '%s' and local copy returned '%s'",
				err.Error(),
				localErr.Error())
		}
	} else {
		// Response received
		mapReader = response.Body
	}

	versionMapBytes, err := ioutil.ReadAll(mapReader)
	if err != nil {
		return err
	}

	var versionMaps VersionMaps
	err = json.Unmarshal(versionMapBytes, &versionMaps)
	if err != nil {
		return err
	}
	defer mapReader.Close()
	updater.versionMaps = versionMaps

	// Write a local cache for the versionmap
	err = ioutil.WriteFile(filepath.Join(
		updater.installPath,
		"versionmap.json"), versionMapBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

// GetLatestVersion returns the latest version installed
func (updater *UT4Updater) GetLatestVersion() (UT4Version, error) {
	versions, err := updater.GetVersionList()
	if err != nil {
		return UT4Version{}, err
	}

	if len(versions) == 0 {
		return UT4Version{}, errors.New("No Unreal Tournament versions installed")
	}

	return versions[0], nil
}

// GetVersionList returns the available installed versions as [version][path]
func (updater *UT4Updater) GetVersionList() ([]UT4Version, error) {
	fileInfo, err := os.Stat(updater.installPath)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() == false {
		return nil, errors.New("The install path must be a directory")
	}

	files, err := ioutil.ReadDir(updater.installPath)
	if err != nil {
		return nil, err
	}

	var versions []UT4Version
	for _, file := range files {
		if file.IsDir() {
			version := UT4Version{
				Path: filepath.Join(updater.installPath, file.Name()),
				VersionMap: updater.versionMaps.GetVersionMapByVersionNumber(
					file.Name()),
			}
			versions = append(versions, version)
		}
	}
	// Reverse the order so that the latest is at the top
	sort.Sort(ByVersion(versions))
	return versions, nil
}

// CheckForUpdate checks if an update is available
func (updater *UT4Updater) CheckForUpdate() (bool, error) {
	latestVersion, err := updater.GetLatestVersion()
	if err != nil {
		return false, err
	}
	osDistribution := OSDistribution{
		Distribution:           "Optout",
		DistributionID:         "optout",
		DistributionPrettyName: "Optout",
		KernelVersion:          "Linux Optout",
		DistributionVersion:    "0.0",
	}
	var versions []string
	if updater.sendStats {
		osDistribution = updater.GetOSDistribution()
		installedVersions, err := updater.GetVersionList()
		if err == nil {
			for _, version := range installedVersions {
				versions = append(versions, version.Version)
			}
		}
	}

	updateCheckRequest := UpdateCheckRequest{
		ClientID:       updater.clientID,
		OS:             osDistribution,
		Versions:       versions,
		CurrentVersion: latestVersion.Version,
	}
	checkJSON, err := json.Marshal(updateCheckRequest)
	if err != nil {
		return false, err
	}

	client := pester.New()
	client.Concurrency = 1
	client.MaxRetries = 3
	client.Backoff = pester.DefaultBackoff
	fmt.Println(fmt.Sprintf("%s/%s/%s", updater.updateURL, "update", "ut4-check"))
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/%s/%s", updater.updateURL, "update", "ut4-check"),
		bytes.NewReader(checkJSON))
	if err != nil {
		return false, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	log.Printf("UpdateStatus %s", resp.Status)
	content, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(content))

	return true, nil
}

// Update creates a backup and the current game, determines the files to be
// updated, downloads the files and applies the updates. Returns the new latest
// version.
// This is safe to run in a goroutine.
func (updater *UT4Updater) Update() (UT4Version, error) {

	// Generate file list with hashes
	// Generate JSON  file list with hashes
	// Generate SHA256 hash for the filelist
	// Generate update manifest
	// Submit the update manifest
	// Clone the current latest version with the new latest version name
	// Download the updated files (wait for package)
	// Apply (remove, add, update) the update
	//
	return UT4Version{}, nil
}

// GetOSDistribution retrieves the kernel and distribution versions
func (updater *UT4Updater) GetOSDistribution() OSDistribution {
	var osDistribution OSDistribution

	// /etc/os-release is the preferred way to check for distribution,
	// if it exists, we'll use it, otherwise just check for another *-release
	// file and user a part of it. This isn't critical to the updater.
	hasReleaseFile := true
	releaseBytes, err := ioutil.ReadFile("/etc/os-release")
	if err != nil {
		// File doesn't exist, check if the next one does
		releaseBytes, err = ioutil.ReadFile("/usr/lib/os-release")
		if err != nil {
			// still no release file
			hasReleaseFile = false
		}
	}
	releaseContents := make(map[string]string)
	if hasReleaseFile {
		for _, line := range strings.Split(string(releaseBytes), "\n") {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.Replace(strings.TrimSpace(parts[1]), "\"", "", -1)
				releaseContents[key] = value
			}
		}
	} else {
		_ = filepath.Walk("/etc",
			func(path string, f os.FileInfo, _ error) error {
				if !f.IsDir() {
					r, walkErr := regexp.MatchString("release", f.Name())
					if walkErr == nil && r {
						// This is a lazy way since this is not really important
						parts := strings.Split(f.Name(), "-")
						if len(parts) == 2 {
							releaseContents["ID"] = strings.Title(parts[0])
							releaseContents["NAME"] = releaseContents["ID"] + " Linux"
							releaseContents["PRETTY_NAME"] = releaseContents["NAME"]
						}
					}
				}
				return nil
			})
		if len(releaseContents) == 0 {
			releaseContents["ID"] = "Generic"
			releaseContents["NAME"] = "Generic Linux"
			releaseContents["PRETTY_NAME"] = "Generic Linux"
		}
	}

	if _, ok := releaseContents["NAME"]; ok {
		osDistribution.Distribution = releaseContents["NAME"]
	}
	if _, ok := releaseContents["ID"]; ok {
		osDistribution.DistributionID = releaseContents["ID"]
	}
	if _, ok := releaseContents["VERSION_ID"]; ok {
		osDistribution.DistributionVersion = releaseContents["VERSION_ID"]
	}
	if _, ok := releaseContents["PRETTY_NAME="]; ok {
		osDistribution.DistributionPrettyName = releaseContents["PRETTY_NAME="]
	}

	out, err := exec.Command("uname", "-r").Output()
	if err != nil {
		// Could not execute uname -r
		osDistribution.KernelVersion = "Unknown"
	} else {
		rawVersion := string(out)
		parts := strings.Split(rawVersion, "-")
		if len(parts) > 0 {
			osDistribution.KernelVersion = parts[0]
		} else {
			osDistribution.KernelVersion = rawVersion
		}
	}

	return osDistribution
}
