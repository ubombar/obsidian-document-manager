package fileapi

import (
	"regexp"

	"github.com/ubombar/obsidian-document-manager/internal/ddm/matcher"
)

type TextFile struct {
	parent *File
}

func NewTextFile(f *File) *TextFile {
	return &TextFile{
		parent: f,
	}
}

// Note that if the file does not exist then it returns an empty array instead of
// an error.
func (m TextFile) Match(regex *regexp.Regexp) (*matcher.MatchResult, error) {
	// Check if the file exists
	if exists, err := m.parent.Exists(); err != nil {
		return nil, err
	} else if !exists {
		return &matcher.MatchResult{}, nil
	}

	// For now get the contents of the file instead of chunking since I don't
	// think memory would be an issue at this moment.
	data, err := m.parent.ReadAll()

	if err != nil {
		return nil, err
	}

	matchesInt := regex.FindAllIndex(data, -1)

	return matcher.NewResult(matchesInt, data)
}
