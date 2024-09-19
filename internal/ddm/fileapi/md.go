package fileapi

import (
	"io"
	"os"
	"regexp"

	"github.com/ubombar/obsidian-document-manager/internal/ddm/matcher"
)

type Markdown struct {
	parent *File
}

func NewMarkdown(f *File) *Markdown {
	return &Markdown{
		parent: f,
	}
}

// Note that if the file does not exist then it returns an empty array instead of
// an error.
func (m Markdown) Match(regex *regexp.Regexp) (*matcher.MatchResult, error) {
	// Check if the file exists
	if exists, err := m.parent.Exists(); err != nil {
		return nil, err
	} else if !exists {
		return &matcher.MatchResult{}, nil
	}

	if f, err := os.Open(m.parent.absPath); err == nil {
		defer f.Close()

		// The text files are not expected to be large so this is ok.
		data, err := io.ReadAll(f)

		if err != nil {
			return nil, err
		}

		matchesInt := regex.FindAllIndex(data, -1)

		return matcher.NewResult(matchesInt, data)
	} else {
		return nil, err
	}
}
