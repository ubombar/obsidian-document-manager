package odm

type Group interface {
	// Text file implements Modifiable
	GroupFilterer

	// Text file implements Matchable
	SetMatchable

	Data() Data
}

type group struct {
	// Set is a set.
	Group

	// The collection
	collection map[string]Set

	// This is the version of the buffer.
	version int
}

func NewEmptyGroup() Group {
	return &group{
		collection: make(map[string]Set),
		version:    0,
	}
}
