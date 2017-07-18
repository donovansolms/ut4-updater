package ut4updater_test

import (
	"testing"

	"github.com/donovansolms/ut4-updater/src/ut4updater"
)

func TestNew(t *testing.T) {
	updater, err := ut4updater.New("")
	if err != nil {
		t.Error(err.Error())
	}
	if updater == nil {
		t.Error("Updater is nil without error")
	}
}
