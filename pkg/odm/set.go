package odm

import (
	"errors"
	"regexp"
)

type Set struct {
	// Text file implements Modifiable
	IndexModifiable

	// Text file implements Matchable
	IndexMatchable

	// This represents the buffer, the data.
	data Data

	// This is the version of the buffer.
	version int
}

func NewEmptySet() (*Set, error) {
	return &Set{
		data:    new([]byte),
		version: 0,
	}, nil
}

func (m *Set) Match(regex *regexp.Regexp) (*[]Match, error) {
	matchesInt := regex.FindAllIndex(*m.data, -1)
	return FromIntArray(&matchesInt)
}

// Assumed the matches are mutually exclusive. This can be done very efifciently using go routings
// But YOLO
func (m *Set) Replace(mm *[]Match, f IndexMatchMapper) (bool, error) {
	return m.perform(mm, f, mode_replace)
}

func (m *Set) InsertBefore(mm *[]Match, f IndexMatchMapper) (bool, error) {
	return m.perform(mm, f, mode_insert_before)
}

func (m *Set) InsertAfter(mm *[]Match, f IndexMatchMapper) (bool, error) {
	return m.perform(mm, f, mode_insert_after)
}

func (m *Set) Remove(mm *[]Match) (bool, error) {
	// Here the function can be nil since it is not called.
	return m.perform(mm, nil, mode_remove)
}

func (m *Set) GetBuffer() Data {
	return m.data
}

const (
	mode_replace       int = 1
	mode_insert_before int = 2
	mode_insert_after  int = 3
	mode_remove        int = 4
)

func (m *Set) perform(mm *[]Match, f IndexMatchMapper, mode int) (bool, error) {
	// Check if there are matches currently
	if m == nil {
		return false, errors.New("given matches array is nil")
	}

	// allocate at least this size, no, this can be optimized later.
	modBuffer := make([]byte, 0)

	// Order the matches TODO

	pointer := 0
	offset := 0

	// Start with the matches
	for _, match := range *mm {
		// copy the before unmatched section
		if pointer < match.Begin {
			modBuffer = append(modBuffer, (*m.data)[pointer:match.Begin]...)
		}

		// If the mode is not remove perform this action
		if mode != mode_remove {
			// Get the matched section
			matchedSection := (*m.data)[match.Begin:match.End]

			// Fucntionate the current section
			newData, ok := f(match, m.data)

			if ok {
				// If the mode is insert after do this
				if mode == mode_insert_after {
					// Increment the offset
					offset += len(matchedSection)

					// Put the original match
					modBuffer = append(modBuffer, matchedSection...)

					// If ok, append the new data
					modBuffer = append(modBuffer, newData...)
				} else if mode == mode_insert_before {
					// Increment the offset
					offset += len(matchedSection)

					// If ok, append the new data
					modBuffer = append(modBuffer, newData...)

					// Put the original match
					modBuffer = append(modBuffer, matchedSection...)
				} else {
					modBuffer = append(modBuffer, newData...)
				}

			} else {
				// What to do if this function is not happy? Copy the original data
				modBuffer = append(modBuffer, (*m.data)[match.Begin:match.End]...)
			}
		}

		// if the mode is remove then simpy skip over the pointer.
		// Then update the pointer
		pointer = match.End
	}

	// copy the remeaning part
	modBuffer = append(modBuffer, (*m.data)[pointer:]...)

	m.data = &modBuffer

	return true, nil
}
