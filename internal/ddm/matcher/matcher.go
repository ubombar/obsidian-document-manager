package matcher

import (
	"errors"
)

// This struct is used for storing the
type MatchResult struct {
}

func NewResult(matches [][]int, data []byte) (*MatchResult, error) {
	return nil, nil
}

// This represents a match
type Match struct {
	Begin int
	End   int
}

func fromIntArray(a [][]int) ([]Match, error) {
	res := make([]Match, len(a))

	for i := 0; i < len(a); i++ {
		if len(a[i]) != 2 {
			return nil, errors.New("given array does not satisfy shape (N, 2)")
		}
		res[i] = Match{
			Begin: a[i][0],
			End:   a[i][1],
		}
	}

	return res, nil
}
