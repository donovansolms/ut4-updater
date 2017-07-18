package ut4updater

import (
	"errors"
	"io/ioutil"
	"os"
)

const (
	runVersionLatest = "latest"
)

// UT4Updater is the main executor for the updater
type UT4Updater struct {
	InstallPath  string
	KeepVersions uint
	RunVersion   string
	SendStats    bool
}

/*
// ConfigureFromStream configures the updater from an io stream. Must be YAML.
func (updater *UT4Updater) ConfigureFromStream(configStream io.Reader) error {
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(configStream)
	if err != nil {
		return err
	}

	updater.keepVersions = uint(viper.GetInt("Versioning.Keep"))
	if updater.keepVersions == 0 {
		return errors.New("Versioning.Keep must be larger than zero")
	}
	updater.runVersion = viper.GetString("Versioning.Run")
	if updater.runVersion == "" {
		log.Println("Versioning.Run is not set, set to `latest`")
	}
	updater.sendStats = viper.GetBool("SendStats")
	if !updater.sendStats {
		log.Println("Stats will not be sent")
	}
	updater.installPath = viper.GetString("InstallPath")
	if updater.installPath == "" {
		return errors.New("InstallPath cannot be blank")
	}

	return nil
}
*/

// GetVersionList returns the available installed versions as [version][path]
func (updater *UT4Updater) GetVersionList() ([]string, error) {
	fileInfo, err := os.Stat(updater.InstallPath)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() == false {
		return nil, errors.New("The install path must be a directory")
	}

	files, err := ioutil.ReadDir(updater.InstallPath)
	if err != nil {
		return nil, err
	}

	var versions []string
	for _, file := range files {
		if file.IsDir() {
			versions = append(versions, file.Name())
		}
	}
	return versions, nil
}

//func (updater *UT4Updater) GetVersion(version latest)
