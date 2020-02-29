package main

/*
Types and functions used in tests.
*/

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/enr/clui"
	"github.com/enr/go-commons/lang"
	"github.com/enr/go-files/files"
)

var (
	testdataPath string
	showOutput   bool
)

func init() {
	testdataPath, _ = filepath.Abs("../../testdata")
	// go test -v ./... -verbose
	// flag.BoolVar(&showOutput, "showoutput", false, "shows output during tests (needs -v flag)")
	// flag.Parse()
	// verbosityLevel := clui.VerbosityLevelMute
	// if showOutput {
	verbosityLevel := clui.VerbosityLevelHigh
	// }
	ui, _ = clui.NewClui(func(ui *clui.Clui) {
		ui.VerbosityLevel = verbosityLevel
	})
}

type directoryItem struct {
	FullPath string
	Hash     string
	Size     int64
}

type directoryListing map[string]directoryItem

func readDirectory(dirpath string) (directoryListing, error) {
	return readDirectoryExcluding(dirpath, []string{})
}
func readDirectoryExcluding(dirpath string, excludes []string) (directoryListing, error) {
	contents := make(map[string]directoryItem)
	source, _ := normalizePath(dirpath)
	if !files.IsDir(source) {
		return contents, fmt.Errorf("%s not a directory", dirpath)
	}

	err := filepath.Walk(source, func(fpath string, f os.FileInfo, err error) error {
		if !files.IsRegular(fpath) {
			return nil
		}
		slashSource := filepath.ToSlash(source)
		fileID := strings.TrimLeft(strings.Replace(filepath.ToSlash(fpath), slashSource, "", 1), "/")
		if lang.SliceContainsString(excludes, fileID) {
			return nil
		}
		fullPath, _ := normalizePath(fpath)
		sha1sum, err := files.Sha1Sum(fullPath)
		if err != nil {
			return err
		}
		contents[fileID] = directoryItem{
			FullPath: fullPath,
			Hash:     sha1sum,
			Size:     f.Size(),
		}
		return nil
	})
	return contents, err
}

func directoriesAreSynchronized(dir1 string, dir2 string, excludes []string) (bool, error) {
	c1, err := readDirectoryExcluding(dir1, excludes)
	if err != nil {
		return false, err
	}
	c2, err := readDirectoryExcluding(dir2, excludes)
	if err != nil {
		return false, err
	}
	if len(c1) != len(c2) {
		return false, fmt.Errorf("%s contains #%d files, %s #%d", dir1, len(c1), dir2, len(c2))
	}
	for key, value := range c1 {
		f2 := c2[key]
		if value.Hash != f2.Hash {
			return false, fmt.Errorf("%s hash: %s != %s", key, value.Hash, f2.Hash)
		}
		if value.Size != f2.Size {
			return false, fmt.Errorf("%s size: %d != %d", key, value.Size, f2.Size)
		}
	}
	return true, nil
}

func tempDirectory() (string, error) {
	return createTempDirectory("resync-test-")
}

func assertStringContains(t *testing.T, s string, substr string) {
	if substr != "" && !strings.Contains(s, substr) {
		t.Fatalf("expected output\n%s\n  does not contain\n%s\n", s, substr)
	}
}

func writeFileInDir(t *testing.T, dirpath string, filename string) {
	f := path.Join(dirpath, filename)
	os.MkdirAll(path.Dir(f), 0777)
	contents := []byte(filename + "\n")
	err := ioutil.WriteFile(f, contents, 0644)
	if err != nil {
		t.Fatalf("error writing %s: %v", f, err)
	}
}

func createTempDirectory(prefix string) (string, error) {
	dirpath, err := ioutil.TempDir(os.TempDir(), prefix)
	if err != nil {
		return "", err
	}
	thepath, err := normalizePath(dirpath)
	return thepath, err
}
