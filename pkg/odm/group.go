package odm

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/ubombar/obsidian-document-manager/pkg/odm/api"
)

type group struct {
	// group is a Group.
	api.Group

	// The collection
	collection map[string]api.Set

	// This is the version of the buffer.
	version int
}

func NewEmptyGroup() api.Group {
	return &group{
		collection: make(map[string]api.Set),
		version:    0,
	}
}

// If already exists, simply returns false, nil
func (g *group) Add(s api.Set) (bool, error) {
	// Check if the attribute
	if s == nil || s.Attributes() == nil {
		return false, errors.New("either s or s.Attributes is nil")
	}
	// If it already contains set, well return false
	if _, ok := g.collection[s.Attributes().Name()]; ok {
		return false, nil
	}

	// Add the s to the collection
	g.collection[s.Attributes().Name()] = s

	return true, nil
}

func (g *group) Filter(f api.GroupRemoveCallback) (int, error) {
	newCollection := make(map[string]api.Set)

	for _, set := range g.collection {
		if ok, err := f(set); !ok {
			continue
		} else if err != nil {
			return 0, err
		}
	}

	g.collection = newCollection

	return len(g.collection), nil
}
func (g *group) Merge(f api.GroupMergeCallback) (int, error) {
	mergedCollection := make(map[string]api.Set)
	firstIndex := 0

	for _, firstSet := range g.collection {
		secondIndex := 0
		mergedData := new([]byte)
		for _, secondSet := range g.collection {
			if firstIndex > secondIndex {
				break
			}

			merge, err := f(firstSet, secondSet)

			if err != nil {
				return 0, err
			}

			if merge {
				*mergedData = append(*mergedData, *secondSet.Data()...)
			}
			secondIndex += 1
		}

		fmt.Printf("(*mergedData): %v\n", (*mergedData))
		mergedSet, err := NewSetFromReader(bytes.NewReader(*mergedData))

		if err != nil {
			return 0, err
		}

		mergedCollection[mergedSet.Attributes().Name()] = mergedSet

		firstIndex += 1
	}

	// Set the group to use the new collection.
	g.collection = mergedCollection

	return len(mergedCollection), nil
}

func (g *group) Sets() []api.Set {
	set := make([]api.Set, len(g.collection))
	i := 0

	for _, s := range g.collection {
		set[i] = s
		i += 1
	}

	return set
}

func (g *group) ForEach(f api.GroupForEachCallback) error {
	newCollection := make(map[string]api.Set)

	for setName, set := range g.collection {
		if newSet, err := f(set); err != nil {
			return err
		} else {
			newCollection[setName] = newSet
		}
	}

	// Set the group to use the new collection.
	g.collection = newCollection

	return nil
}
