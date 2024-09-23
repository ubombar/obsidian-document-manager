package file

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const DEFAULT_FILE_PERM int32 = 0644

type File struct {
	// Standard library file interface.
	io.ReadWriteCloser

	// This file is also implements stringer interface.
	fmt.Stringer

	// Absolute file path.
	absPath string

	// Actual std file object
	file *os.File
}

// Create a new file. Any path can be given as the path, it will be converted to the abs
// path. So the current working directory of the program must not be changed Otherwise, it
// can lead to undesired behavior.
func NewFile(path string) (*File, error) {
	absPath, err := cleanedAbsPath(path)

	if err != nil {
		return nil, err
	}

	return &File{
		absPath: absPath,
		file:    nil, // default values are nil since we don't know if this file exists.
	}, nil
}

func cleanedAbsPath(path string) (string, error) {
	// At this point the cleanedPath is not guarenteed to be abs path, so join with pwd.
	cleanedPath := filepath.Clean(strings.TrimLeft(path, " "))

	// A file cannot end with / so die.
	if strings.HasSuffix(cleanedPath, "/") {
		return "", errors.New("a file cannot end with suffix '/'")
	}

	// If it is an absolute path.
	if strings.HasPrefix(cleanedPath, "/") {
		return cleanedPath, nil
	}

	// Join with pwd, so the file would be holding the abs path.
	if pwd, err := os.Getwd(); err == nil {
		cleanedAbsPath := filepath.Join(pwd, cleanedPath)

		return cleanedAbsPath, nil
	} else {
		return "", err
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

// Check if the file exists on the disk
func (f *File) Exists() (bool, error) {
	if _, err := os.Stat(f.absPath); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}

// This assumes the file is already open.
func (f *File) IsOpen() bool {
	return f.file != nil
}

// Opens the actual file, or creates it if not exist
func (f *File) Open() error {
	if file, err := os.OpenFile(f.absPath, os.O_RDWR|os.O_CREATE, fs.FileMode(DEFAULT_FILE_PERM)); err == nil {
		f.file = file
		return nil
	} else {
		return err
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

func (f *File) Close() error {
	if f.IsOpen() {
		return errors.New("file is not opened")
	} else {
		if err := f.file.Close(); err != nil {
			return err
		}
	}
	f.file = nil
	return nil
}

// If file doesn't exist, then read empty bytes
func (f *File) Read(p []byte) (n int, err error) {
	if !f.IsOpen() {
		return 0, errors.New("file is not open")
	}

	return f.file.Read(p)
}

func (f *File) Write(p []byte) (n int, err error) {
	if !f.IsOpen() {
		return 0, errors.New("file is not open")
	}

	return f.file.Write(p)
}
