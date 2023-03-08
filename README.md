# scaffolder

A Go library for scaffolding.

Built on [hackpadfs's](https://github.com/hack-pad/hackpadfs) [FS](https://pkg.go.dev/io/fs#FS) implementations, scaffolder adds higher level populate and copy capabilities.

With scaffolder, you can do things like read in a file system from disk, manipulate the file system, add additional directories and files, then write the file system to another disk-based, in-memory, or other types of file systems.


## Usage Example

```go
// create a disk-based file system rooted at output path (remove trailing "/")
osFS := os.NewFS()
destFS, err := osFS.Sub(outDir[1:])
check(err)

// define a hierarchical directory structure
contents := scaffolder.FSContents{
  "file.txt":           []byte("file1"),
  "dir1/file.txt":      []byte("file2"),
  "dir1/dir2/file.txt": []byte("file3"),
}

// populate the file system with contents
err = scaffolder.PopulateFS(destFS, contents)
check(err)

// inspect the fs
err = scaffolder.InspectFS(destFS, log, false)
check(err)

// copy fs to another place on disk
destFS2, err := hackpadfs.Sub(destFS, "copy")
check(err)
err = scaffolder.CopyFS(destFS, destFS2)
check(err)
```


## Development

```
 Choose a make command to run

  vet           vet code
  test          run unit tests
  build         build a binary
  autobuild     auto build when source files change
  dockerbuild   build project into a docker container image
  start         build and run local project
  deploy        build code into a container and deploy it to the cloud dev environment
```
