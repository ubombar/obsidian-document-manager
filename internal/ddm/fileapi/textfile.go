package fileapi

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/ubombar/obsidian-document-manager/internal/ddm/engine"
)

type TextFile struct {
	// Text file implements Modifiable
	engine.Modifiable

	// Text file implements Matchable
	engine.Matchable

	// Ref to the parent file
	parent *File

	// This represents the buffer.
	buffer *[]byte

	// is buffer changed
	dirty bool
}

func NewTextFile(f *File) (*TextFile, error) {
	file := &TextFile{
		parent: f,
		buffer: new([]byte),
		dirty:  false,
	}

	// Load the contents of the file into the buffer
	if exists, err := f.Exists(); err == nil && exists {
		file.buffer, err = f.ReadAll()

		if err != nil {
			return nil, err
		}

		return file, err
	} else if err != nil {
		return nil, err
	} else {
		return file, nil
	}
}

// Loads the contents of the file to the buffer, if the file doesn't exist
// then set the buffer empty.
func (m *TextFile) LoadBuffer() error {
	var err error
	if m.buffer, err = m.parent.ReadAll(); err == nil {
		return err
	}
	return nil
}

// Flus the contents of the buffer to the disk.
func (m *TextFile) WriteBuffer() error {
	return m.parent.WriteAll(m.buffer)
}

func (m *TextFile) Match(regex *regexp.Regexp) (*[]engine.Match, error) {
	matchesInt := regex.FindAllIndex(*m.buffer, -1)
	return engine.FromIntArray(&matchesInt)
}

// Assumed the matches are mutually exclusive. This can be done very efifciently using go routings
// But YOLO
func (m *TextFile) Replace(mm *[]engine.Match, f engine.MatchMapper) (bool, error) {
	return m.perform(mm, f, mode_replace)
}

func (m *TextFile) InsertBefore(mm *[]engine.Match, f engine.MatchMapper) (bool, error) {
	return m.perform(mm, f, mode_insert_before)
}

func (m *TextFile) InsertAfter(mm *[]engine.Match, f engine.MatchMapper) (bool, error) {
	return m.perform(mm, f, mode_insert_after)
}

func (m *TextFile) Remove(mm *[]engine.Match) (bool, error) {
	return m.perform(mm, func(md engine.Match, buffer *[]byte) ([]byte, bool) { return []byte(""), true }, mode_remove)
}

func (m *TextFile) GetBuffer() *[]byte {
	return m.buffer
}

const (
	mode_replace       int = 1
	mode_insert_before int = 2
	mode_insert_after  int = 3
	mode_remove        int = 4
)

func (m *TextFile) perform(mm *[]engine.Match, f engine.MatchMapper, mode int) (bool, error) {
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
			modBuffer = append(modBuffer, (*m.buffer)[pointer:match.Begin]...)
		}

		// If the mode is not remove perform this action
		if mode != mode_remove {
			// Get the matched section
			matchedSection := (*m.buffer)[match.Begin:match.End]

			// Fucntionate the current section
			newData, ok := f(match, m.buffer)

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
				modBuffer = append(modBuffer, (*m.buffer)[match.Begin:match.End]...)
			}
		}

		// if the mode is remove then simpy skip over the pointer.
		// Then update the pointer
		pointer = match.End + offset
	}

	// copy the remeaning part
	modBuffer = append(modBuffer, (*m.buffer)[pointer:]...)
	fmt.Printf("modbuffer:\n%v\n", string(modBuffer))

	m.buffer = &modBuffer

	return true, nil
}
