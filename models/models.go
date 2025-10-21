package models

type Config struct {
	SrcPath  string
	DestPath string

	ShouldServe     bool
	ShouldSeedFiles bool
}
