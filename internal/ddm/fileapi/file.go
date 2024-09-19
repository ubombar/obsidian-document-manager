package fileapi

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const DEFAULT_FILE_PERM int32 = 0644

type File struct {
	absPath string

	// This property is used to manipulate the contents. This can be markdown, pdf etc.
	fileInterface interface{}
}

// Create a new file. Any path can be given as the path, it will be converted to the abs
// path. So the current working directory of the program must not be changed Otherwise, it
// can lead to undesired behavior.
func NewFile(path string) (*File, error) {
	// At this point the cleanedPath is not guarenteed to be abs path, so join with pwd.
	cleanedPath := filepath.Clean(strings.TrimLeft(path, " "))

	// A file cannot end with / so die.
	if strings.HasSuffix(cleanedPath, "/") {
		return nil, errors.New("a file cannot end with suffix '/'")
	}

	// If it is an absolute path.
	if strings.HasPrefix(cleanedPath, "/") {
		return &File{
			absPath: cleanedPath,
		}, nil
	}

	// Join with pwd, so the file would be holding the abs path.
	if pwd, err := os.Getwd(); err == nil {
		cleanedAbsPath := filepath.Join(pwd, cleanedPath)

		return &File{
			absPath: cleanedAbsPath,
		}, nil
	} else {
		return nil, err
	}

}

// Get the base name including the extension.
func (f *File) BaseName() string {
	return filepath.Base(f.absPath)
}

// Gets the extension, after the dot '.'
func (f *File) Extension() string {
	return filepath.Ext(f.absPath)
}

// Gets the parent file of the file. If it is the root file
// then return itself instead of returning an error.
func (f *File) Parent() (*Folder, error) {
	return f.ParentTimes(1)
}

// Gets the n times parent file of the current file.
func (f *File) ParentTimes(n int) (*Folder, error) {
	current := f.absPath
	for i := 0; i < n; i++ {
		current = filepath.Dir(current)
	}
	return NewFolder(current)
}

// ToString method.
func (f *File) String() string {
	return f.absPath
}

// Check if the file exists
func (f *File) Exists() (bool, error) {
	if _, err := os.Stat(f.absPath); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}

// Creates the actual file. Returns true if new files are created.
// If it already exists, returns false, or error.
func (f *File) Create() (bool, error) {
	// If there is an error or file exists then return false and the error
	if ok, err := f.Exists(); err != nil || ok {
		return false, err
	}

	// If not then attempt to create the file.
	if fi, err := os.Create(f.absPath); err == nil {
		defer fi.Close()
		if err = fi.Chmod(fs.FileMode(DEFAULT_FILE_PERM)); err != nil {
			return false, err
		}
		return true, nil
	} else if errors.Is(err, os.ErrExist) {
		return false, nil
	} else {
		return false, err
	}
}

// Set the permissions of the file, default is 644.
func (f *File) SetPermissions(rwx int32) (bool, error) {
	if ok, err := f.Exists(); err != nil || !ok {
		return false, err
	}

	if fi, err := os.Open(f.absPath); err == nil {
		defer fi.Close()
		if err = fi.Chmod(fs.FileMode(rwx)); err != nil {
			return false, err
		}
		return true, nil
	} else {
		return false, err
	}
}

// Returns true if the file is a markdown file.
func (f *File) IsMarkdown() bool {
	_, ok := f.fileInterface.(*Markdown)
	return ok
}

// Returns false if it is already set to markdown
func (f *File) SetAsMarkdown() bool {
	if f.fileInterface != nil {
		return false
	}

	f.fileInterface = NewMarkdown(f)
	return true
}

func (f *File) Markdown() *Markdown {
	if m, ok := f.fileInterface.(*Markdown); ok {
		return m
	}
	return nil
}
