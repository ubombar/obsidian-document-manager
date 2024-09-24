package inmemory

import (
	"fmt"
	"time"

	"github.com/ubombar/obsidian-document-manager/pkg/odm/api"
)

type inMemorySetAttributes struct {
	api.SetAttirbuter
	created time.Time
	updated time.Time
	version int
}

func NewInMemorySetAttriubtes() api.SetAttirbuter {
	return &inMemorySetAttributes{
		created: time.Now(),
		updated: time.Now(),
		version: 0,
	}
}

// On a file based one, this returns the filepath.
func (sa inMemorySetAttributes) Name() string {
	return fmt.Sprintf("%v", &sa)
}

// Returns the creation timestamp
func (sa inMemorySetAttributes) Created() (time.Time, error) {
	return time.Now(), nil
}

// Returns the updated timestamp.
func (sa inMemorySetAttributes) Updated() (time.Time, error) {
	return sa.updated, nil
}

// Returns the updated timestamp.
func (sa *inMemorySetAttributes) Update(t time.Time) error {
	sa.updated = time.Now()
	return nil
}

// Returns the version.
func (sa inMemorySetAttributes) Version() int {
	return sa.version
}

func (sa *inMemorySetAttributes) IncrementVersion() {
	sa.version += 1
}
