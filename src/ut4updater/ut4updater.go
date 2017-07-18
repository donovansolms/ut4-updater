package ut4updater

// UT4Updater is the main executor for the updater
type UT4Updater struct {
}

// New creates and initializes a new UT4Updater instance
func New(configPath string) (*UT4Updater, error) {
	updater := UT4Updater{}

	return &updater, nil
}
