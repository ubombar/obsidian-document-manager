package odm

import (
	"errors"
	"io"
	"regexp"
	"time"

	"github.com/ubombar/obsidian-document-manager/pkg/odm/api"
	"github.com/ubombar/obsidian-document-manager/pkg/odm/attributes/disk"
	"github.com/ubombar/obsidian-document-manager/pkg/odm/attributes/inmemory"
	"github.com/ubombar/obsidian-document-manager/pkg/odm/file"
)

type set struct {
	// Set is a set.
	api.Set

	// This is used to receive the set attributes.
	attributes api.SetAttirbuter

	// This represents the buffer, the data.
	data api.Data
}

func NewEmptySet() api.Set {
	return &set{
		data:       new([]byte),
		attributes: inmemory.NewInMemorySetAttriubtes(),
	}
}

func NewClonedSet(s api.Set) api.Set {
	return &set{
		data:       s.Data(),
		attributes: s.Attributes(),
	}
}

func NewSetFromReader(reader io.Reader) (api.Set, error) {
	if data, err := io.ReadAll(reader); err != nil {
		return nil, err
	} else {
		return &set{
			data:       &data,
			attributes: inmemory.NewInMemorySetAttriubtes(),
		}, nil
	}
}

func NewSetFromFile(file *file.File) (api.Set, error) {
	if !file.IsOpen() {
		file.Open()
		defer file.Close()
	}
	if data, err := io.ReadAll(file); err != nil {
		return nil, err
	} else {
		attrib, err := disk.NewDiskSetAttributes(file)

		if err != nil {
			return nil, err
		}

		return &set{
			data:       &data,
			attributes: attrib,
		}, nil
	}
}

func NewSetFromFileOrEmpty(file *file.File) api.Set {
	if data, err := io.ReadAll(file); err != nil {
		return NewEmptySet()
	} else {
		attrib, err := disk.NewDiskSetAttributes(file)

		if err != nil {
			return NewEmptySet()
		}

		return &set{
			data:       &data,
			attributes: attrib,
		}
	}
}

func (m *set) CompiledMatch(regex string) (*[]api.Match, error) {
	if regex, err := regexp.Compile(regex); err == nil {
		return m.Match(regex)
	} else {
		return nil, err
	}
}

func (m *set) Match(regex *regexp.Regexp) (*[]api.Match, error) {
	matchesInt := regex.FindAllIndex(*m.data, -1)
	return api.FromIntArray(&matchesInt)
}

// Assumed the matches are mutually exclusive. This can be done very efifciently using go routings
// But YOLO
func (m *set) Replace(mm *[]api.Match, f api.SetActionCallback) (bool, error) {
	return m.perform(mm, f, mode_replace)
}

func (m *set) InsertBefore(mm *[]api.Match, f api.SetActionCallback) (bool, error) {
	return m.perform(mm, f, mode_insert_before)
}

func (m *set) InsertAfter(mm *[]api.Match, f api.SetActionCallback) (bool, error) {
	return m.perform(mm, f, mode_insert_after)
}

func (m *set) Remove(mm *[]api.Match) (bool, error) {
	// Here the function can be nil since it is not called.
	return m.perform(mm, nil, mode_remove)
}

func (m *set) Data() api.Data {
	return m.data
}

// Cloning causes the version to reset.
func (m *set) Clone() api.Set {
	return NewClonedSet(m)
}

const (
	mode_replace       int = 1
	mode_insert_before int = 2
	mode_insert_after  int = 3
	mode_remove        int = 4
)

func (m *set) perform(mm *[]api.Match, f api.SetActionCallback, mode int) (bool, error) {
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

	// Set the pointer
	m.data = &modBuffer

	// Bump version
	m.attributes.IncrementVersion()
	m.attributes.Update(time.Now())

	return true, nil
}

func (m *set) Attributes() api.SetAttirbuter {
	return m.attributes
}
