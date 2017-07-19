package ut4updater

// UT4Version holds information about an installed UT4 version
type UT4Version struct {
	VersionMap
	Path string
}

// ByVersion allows for sorting by the build version number
type ByVersion []UT4Version

func (a ByVersion) Len() int           { return len(a) }
func (a ByVersion) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByVersion) Less(i, j int) bool { return a[i].Version > a[j].Version }
