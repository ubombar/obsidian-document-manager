package matcher

import (
	"errors"
)

// This represents a match without the data.
type Match struct {
	Begin int
	End   int
}

// This is a debug type, it contains the Data segment as well. This
// is not very efficient and not used on the action mapping function.
type MatchWithData struct {
	Begin int
	End   int
	Data  []byte
}

func fromMatch(m Match, data []byte) MatchWithData {
	return MatchWithData{
		Begin: m.Begin,
		End:   m.End,
		Data:  data[m.Begin:m.End],
	}
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

// This is the match result. The idea is that it has a certain validity date. If the file
// is modified, then this becomes unvalid. This can be tracked by the valid field.
//
// Note that this does not include the changes from outside. So the name would be
// latest.
type MatchResult struct {
	Modifier

	// Matches
	matches []Match

	// This is here for testing if this match result is the latest one from the file. If any
	// modifications are done in this context (from this program) then the latest is false.
	dirty bool
}

// Check of the Match result implements all the functions. (This doesn't work?)
var _ Modifier = (*MatchResult)(nil)

// Create a new match result.
func NewResult(matches [][]int, data []byte) (*MatchResult, error) {
	formattedMatches, err := fromIntArray(matches)

	if err != nil {
		return nil, err
	}

	return &MatchResult{
		matches: formattedMatches,
		dirty:   false,
	}, nil
}

func (m MatchResult) ToMatches() []Match {
	return m.matches
}

func (m MatchResult) ToMatchesWithData(data []byte) []MatchWithData {
	matches := m.ToMatches()
	matchesWithData := make([]MatchWithData, len(matches))

	for i, match := range m.ToMatches() {
		matchesWithData[i] = fromMatch(match, data)
	}

	return matchesWithData
}

func (m MatchResult) Replace(MatchMapper) (bool, error) {
	return false, nil
}

func (m MatchResult) InsertBefore(MatchMapper) (bool, error) {
	return false, nil
}

func (m MatchResult) InsertAfter(MatchMapper) (bool, error) {
	return false, nil
}

func (m MatchResult) Remove() (bool, error) {
	return false, nil
}
