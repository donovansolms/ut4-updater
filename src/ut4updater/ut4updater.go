package ut4updater

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
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
	versionMaps  VersionMaps
}

// New creates aand initializes a new instance of UT4Updater
func New(installPath string,
	keepVersions uint,
	runVersion string,
	sendStats bool,
	versionMapURL string) (*UT4Updater, error) {
	updater := &UT4Updater{
		installPath:  installPath,
		keepVersions: keepVersions,
		runVersion:   runVersion,
		sendStats:    sendStats,
	}
	fullPath, err := filepath.Abs(updater.installPath)
	if err != nil {
		return updater, err
	}
	updater.installPath = fullPath

	err = updater.updateVersionMap(versionMapURL)
	if err != nil {
		return updater,
			fmt.Errorf("Unable to update version map '%s': %s",
				versionMapURL,
				err.Error())
	}

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
		mapReader, localErr = os.Open(filepath.Join(updater.installPath, "versionmap.json"))
		if localErr != nil {
			return fmt.Errorf("Remote returned '%s' and local copy returned '%s'",
				err.Error(),
				localErr.Error())
		}
	} else {
		// Response received
		mapReader = response.Body
	}

	var versionMaps VersionMaps
	err = json.NewDecoder(mapReader).Decode(&versionMaps)
	if err != nil {
		return err
	}
	defer mapReader.Close()
	updater.versionMaps = versionMaps
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

//func (updater *UT4Updater) GetVersion(version latest)
