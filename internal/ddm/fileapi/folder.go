package fileapi

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Folder struct {
	absPath string
}

// Create a new folder. Any path can be given as the path, it will be converted to the abs
// path. So the current working directory of the program must not be changed Otherwise, it
// can lead to undesired behavior.
func NewFolder(path string) (*Folder, error) {
	// At this point the cleanedPath is not guarenteed to be abs path, so join with pwd.
	cleanedPath := filepath.Clean(strings.TrimLeft(path, " "))

	// If it is an absolute path.
	if strings.HasPrefix(cleanedPath, "/") {
		return &Folder{
			absPath: cleanedPath,
		}, nil
	}

	// Join with pwd, so the folder would be holding the abs path.
	if pwd, err := os.Getwd(); err == nil {
		cleanedAbsPath := filepath.Join(pwd, cleanedPath)

		return &Folder{
			absPath: cleanedAbsPath,
		}, nil
	} else {
		return nil, err
	}

}

func (f *Folder) BaseName() string {
	return filepath.Base(f.absPath)
}

// Gets the parent folder of the folder. If it is the root folder
// then return itself instead of returning an error.
func (f *Folder) Parent() (*Folder, error) {
	return f.ParentTimes(1)
}

// Gets the n times parent folder of the current folder.
func (f *Folder) ParentTimes(n int) (*Folder, error) {
	current := f.absPath
	for i := 0; i < n; i++ {
		current = filepath.Dir(current)
	}
	return NewFolder(current)
}

// Check if it is root.
func (f *Folder) IsRoot() bool {
	return f.absPath == "/"
}

// ToString method.
func (f *Folder) String() string {
	return f.absPath
}

// Check if the folder exists
func (f *Folder) Exists() (bool, error) {
	if _, err := os.Stat(f.absPath); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}

// Creates the actual folder. Returns true if new folders are created.
// If it already exists, returns false, or error.
func (f *Folder) Create() (bool, error) {
	// If there is an error or folder exists then return false and the error
	if ok, err := f.Exists(); err != nil || ok {
		return false, err
	}

	// If not then attempt to create the folder.
	if err := os.MkdirAll(f.absPath, os.ModePerm); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrExist) {
		return false, nil
	} else {
		return false, err
	}
}
