package odm

import "github.com/ubombar/obsidian-document-manager/pkg/odm/api"

type Group interface {
	// Text file implements Modifiable
	api.GroupFilterer

	// Text file implements Matchable
	api.SetMatcher

	Data() api.Data
}

type group struct {
	// Set is a set.
	Group

	// The collection
	collection map[string]api.Set

	// This is the version of the buffer.
	version int
}

func NewEmptyGroup() Group {
	return &group{
		collection: make(map[string]api.Set),
		version:    0,
	}
}
