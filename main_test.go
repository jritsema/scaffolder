package scaffolder

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/hack-pad/hackpadfs"
	"github.com/hack-pad/hackpadfs/mem"
	hackpados "github.com/hack-pad/hackpadfs/os"
)

func TestCopyFromMemory(t *testing.T) {

	contents := FSContents{
		"file.txt":           []byte("file1"),
		"dir1/file.txt":      []byte("file2"),
		"dir1/dir2/file.txt": []byte("file3"),
	}

	//create an in-memory fs destination
	dest, err := mem.NewFS()
	if err != nil {
		t.Error(err)
	}

	//populate the file system with contents
	err = PopulateFS(dest, contents)
	if err != nil {
		t.Error(err)
	}

	err = InspectFS(dest, t.Log, false)
	if err != nil {
		t.Error(err)
	}

	findings := 0

	err = hackpadfs.WalkDir(dest, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if contains(contents, path) {
			findings++
			t.Log("found", path)
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}

	t.Log("expectations:", len(contents))
	t.Log("findings:", findings)

	if findings != len(contents) {
		t.Error(errors.New("path counts don't match. did you add/remove something in the local test directory?"))
	}
}

func TestCopyFromDisk(t *testing.T) {

	//create a src fs pointing at local test directory
	osfs := hackpados.NewFS()
	workingDirectory, _ := os.Getwd()
	src, err := osfs.Sub(path.Join(workingDirectory[1:], "test"))
	if err != nil {
		t.Error(err)
	}

	//create an in-memory fs destination
	dest, err := mem.NewFS()
	if err != nil {
		t.Error(err)
	}

	//copy src fs to dest fs
	err = CopyFS(src, dest)
	if err != nil {
		panic(err)
	}

	err = InspectFS(dest, t.Log, false)
	if err != nil {
		t.Error(err)
	}

	//test that this directory structure is copied
	paths := []string{
		"file.txt",
		"dir1/file.txt",
		"dir2/file.txt",
		"dir2/dir3/file.txt",
	}

	findings := 0

	err = hackpadfs.WalkDir(dest, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if sliceContains(&paths, path) {
			findings++
			t.Log("found", path)
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}

	t.Log("expectations:", len(paths))
	t.Log("findings:", findings)

	if findings != len(paths) {
		t.Error(errors.New("path counts don't match. did you add/remove something in the local test directory?"))
	}
}
