package fileapi

import (
	"errors"
	"io"
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
	// Check the extension, if it is a text or markdown file set the matcher

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

// Do not forget to close the reader. If the file does not exists returns
// an error
func (f *File) Reader() (io.Reader, error) {
	if exists, err := f.Exists(); err != nil {
		return nil, err
	} else if !exists {
		return nil, errors.New("file does not exist")
	}

	if file, err := os.Open(f.absPath); err == nil {
		return file, err
	} else {
		return nil, err
	}
}

func (f *File) ReadAll() (*[]byte, error) {
	if reader, err := f.Reader(); err != nil {
		return nil, err
	} else {
		// Try to close the returned reader if it implements the io.Closer interface.
		if closer, ok := reader.(io.Closer); ok {
			defer closer.Close()
		}

		// Read all the data
		data, err := io.ReadAll(reader)

		if err != nil {
			return nil, err
		}

		return &data, nil
	}
}

func (f *File) WriteAll(data *[]byte) error {
	file, err := os.OpenFile(f.absPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fs.FileMode(DEFAULT_FILE_PERM))
	if err != nil {
		return err
	}
	defer file.Close()

	// Write byte array to file
	_, err = file.Write(*data)
	if err != nil {
		return err
	}

	return nil
}
