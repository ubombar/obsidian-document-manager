package engine

import (
	"errors"
	"fmt"
	"regexp"
)

// This is the map function, it takes a Match and the data segment. This function will
// be invoken when writing down to the file.
// Do not change the contents of the buffer.
type MatchMapper func(md Match, buffer *[]byte) ([]byte, bool)

// Performs a modification action.
type Modifiable interface {
	// Replace the matched section with the given map function
	Replace([]Match, MatchMapper) (bool, error)

	// Insert the matched section before/after with the map function.
	InsertBefore([]Match, MatchMapper) (bool, error)
	InsertAfter([]Match, MatchMapper) (bool, error)

	// Remove that chuck of the mathing
	Remove([]Match) (bool, error)
}

type Match struct {
	Begin int
	End   int
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
func (md Match) Match(buffer *[]byte) string {
	return string((*buffer)[md.Begin:md.End])
}

// The matcher interface
type Matchable interface {
	Match(regex *regexp.Regexp) (*[]Match, error)
}

func Returnf(format string, obj ...any) ([]byte, bool) {
	return []byte(fmt.Sprintf(format, obj...)), true
}
