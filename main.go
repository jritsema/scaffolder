package scaffolder

import (
	"io/fs"
	"path"
	"path/filepath"

	"github.com/hack-pad/hackpadfs"
)

// Logger is a logging function
type Logger = func(args ...interface{})

// FSContents is a directory structure expressed as a map of []byte (content) keyed by path
type FSContents = map[string][]byte

// CreateFile creates a file in the filesystem
func CreateFile(fs hackpadfs.FS, fullPath string, contents []byte) error {

	//create the directory
	dir := filepath.Dir(fullPath)
	err := hackpadfs.MkdirAll(fs, dir, hackpadfs.FileMode.Perm(0755))
	if err != nil {
		return err
	}

	//create the file
	file, err := hackpadfs.Create(fs, fullPath)
	if err != nil {
		return err
	}
	_, err = hackpadfs.WriteFile(file, contents)
	if err != nil {
		return err
	}

	return nil
}

// CreateFile creates a file in the filesystem
func CreateFileWithParts(fs hackpadfs.FS, contents []byte, pathParts ...string) error {
	return CreateFile(fs, path.Join(pathParts...), contents)
}

// PopulateFS populates a file system with contents
func PopulateFS(fs hackpadfs.FS, fsc FSContents) error {
	for file, contents := range fsc {
		err := CreateFile(fs, file, contents)
		if err != nil {
			return err
		}
	}
	return nil
}

// CopyFS copies a file system to another file system
func CopyFS(src, dest hackpadfs.FS) error {
	err := hackpadfs.WalkDir(src, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			b, err := hackpadfs.ReadFile(src, path)
			if err != nil {
				return err
			}
			err = CreateFile(dest, path, b)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// InspectFS inspects a file system, outputting each file using the Logger
func InspectFS(filesystem hackpadfs.FS, log Logger, logContents bool) error {

	err := hackpadfs.WalkDir(filesystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		log()
		log("IsDir =", d.IsDir())
		log("path =", path)
		log("Name =", d.Name())

		//read files
		if !d.IsDir() {

			//open + stat
			f, err := filesystem.Open(path)
			if err != nil {
				return err
			}
			fi, err := f.Stat()
			if err != nil {
				return err
			}
			log("size =", fi.Size())
			log("ModTime =", fi.ModTime())

			b, err := hackpadfs.ReadFile(filesystem, path)
			if err != nil {
				log(err)
				return err
			}
			if logContents {
				log("contents =", string(b))
			}
		}
		return nil
	})

	return err
}

// sliceContains returns true if a slice contains a string
func sliceContains(s *[]string, e string) bool {
	for _, str := range *s {
		if str == e {
			return true
		}
	}
	return false
}

// contains returns true if a slice contains a string
func contains(contents FSContents, e string) bool {
	for c := range contents {
		if c == e {
			return true
		}
	}
	return false
}
