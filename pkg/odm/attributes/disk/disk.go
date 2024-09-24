package disk

import (
	"time"

	"github.com/ubombar/obsidian-document-manager/pkg/odm/api"
	"github.com/ubombar/obsidian-document-manager/pkg/odm/file"
)

type diskSetAttributes struct {
	api.SetAttirbuter
	file    *file.File
	version int
}

// On a file based one, this returns the filepath.
func (sa diskSetAttributes) Name() string {
	return sa.file.String()
}

// Returns the creation timestamp
func (sa diskSetAttributes) Created() (time.Time, error) {
	return sa.file.Created()
}

// Returns the updated timestamp.
func (sa diskSetAttributes) Updated() (time.Time, error) {
	return sa.file.Updated()
}

// Returns the updated timestamp.
func (sa *diskSetAttributes) Update(_ time.Time) error {
	return nil
}

// Returns the version.
func (sa diskSetAttributes) Version() int {
	return sa.version
}

func (sa *diskSetAttributes) IncrementVersion() {
	sa.version += 1
}
