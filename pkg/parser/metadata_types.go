package parser

import (
	"github.com/puzpuzpuz/xsync/v3"
)

type Meta struct {
	// Title
	Title string

	// Sitepath
	Sitepath string

	// Description
	Description string

	// Tags
	Tags []string

	// File
	F *File

	// Frontmatter
	Frontmatter *Frontmatter
}

type Metadata struct {
	Store *xsync.MapOf[string, *Meta]
}
