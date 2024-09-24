package api

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

// This is the main type of data. it is a byte array pointer.
type Data *[]byte

// This is the map function, it takes a Match and the data segment. This function will
// be invoken when writing down to the file.
// Do not change the contents of the buffer.
type SetActionCallback func(md Match, buffer Data) ([]byte, bool)

// Performs a modification action.
type SetModifier interface {
	// Replace the matched section with the given map function
	Replace(*[]Match, SetActionCallback) (bool, error)

	// Insert the matched section before with the map function.
	InsertBefore(*[]Match, SetActionCallback) (bool, error)

	// Insert the matched section after with the map function.
	InsertAfter(*[]Match, SetActionCallback) (bool, error)

	// Remove that chuck of the mathing
	Remove(*[]Match) (bool, error)
}

// The matcher interface
type SetMatcher interface {
	Match(regex *regexp.Regexp) (*[]Match, error)

	CompiledMatch(regex string) (*[]Match, error)
}

type Match struct {
	Begin int
	End   int
}

// The set interface
type Set interface {
	// Text file implements Modifiable
	SetModifier

	// Text file implements Matchable
	SetMatcher

	// Gets the data
	Data() Data

	// Get attributes
	Attributes() SetAttirbuter
}

type SetAttirbuter interface {
	// On a file based one, this returns the filepath.
	Name() string

	// Returns the creation timestamp
	Created() (time.Time, error)

	// Returns the updated timestamp.
	Updated() (time.Time, error)

	// Set the latest updated time.
	Update(t time.Time) error

	// Returns the version.
	Version() int

	// Increment the version
	IncrementVersion()
}

type GroupFilterCallback func() (bool, error)

type GroupFilterer interface {
	Filter(regex *regexp.Regexp)
}

func FromIntArray(a *[][]int) (*[]Match, error) {
	matches := make([]Match, len(*a))

	for i, mint := range *a {
		if len(mint) != 2 {
			return nil, errors.New("given int array is not the shape (N, 2)")
		}
		matches[i] = Match{
			Begin: mint[0],
			End:   mint[1],
		}
	}
	return &matches, nil
}

// Gets the word as string.
func (md Match) Segment(buffer Data) string {
	return string((*buffer)[md.Begin:md.End])
}

func EasyReturn(format string, obj ...any) ([]byte, bool) {
	return []byte(fmt.Sprintf(format, obj...)), true
}
