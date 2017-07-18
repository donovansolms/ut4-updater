package ut4updater_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/donovansolms/ut4-updater/src/ut4updater"
)

var updater ut4updater.UT4Updater

func TestMain(m *testing.M) {
	updater = ut4updater.UT4Updater{
		InstallPath:  "../../test-resources/installs",
		KeepVersions: 2,
		RunVersion:   "latest",
		SendStats:    true,
	}
	m.Run()

}
func TestGetVersionList(t *testing.T) {
	versions, err := updater.GetVersionList()
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(versions)
	sort.Sort(sort.Reverse(sort.StringSlice(versions)))
	fmt.Println(versions)
}
