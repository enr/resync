package main

import (
	"testing"

	"github.com/enr/go-files/files"
)

func TestNoop(t *testing.T) {
	tmp, err := tempDirectory()
	if err != nil {
		t.Errorf(`unexpected error creating temp dir: %v`, err)
	}
	options := syncOptions{
		Source:    testdataPath,
		Target:    tmp,
		Timestamp: "timestamp",
		Noop:      true,
	}
	syncError := syncronizeDirectories(options)
	if syncError != nil {
		t.Errorf(`syncronizeDirectories error: %v`, syncError)
	}
	dircontents, err := readDirectoryExcluding(tmp, []string{})
	if err != nil {
		t.Errorf(`unexpected error reading %s contents: %v`, tmp, err)
	}
	if len(dircontents) != 0 {
		t.Errorf(`unexpected content in %s after "noop" operation`, tmp)
	}
}

func TestBasic(t *testing.T) {
	tmp, err := tempDirectory()
	if err != nil {
		t.Errorf(`unexpected error creating temp dir: %v`, err)
	}
	if !files.IsDir(tmp) {
		t.Errorf(`error creating temp dir: %s`, tmp)
	}
	sd, err := normalizeDirsPath(testdataPath)
	if err != nil {
		t.Errorf(`error managing testdataPath: %v`, err)
	}
	td, err := normalizeDirsPath(tmp)
	if err != nil {
		t.Errorf(`error managing tmp: %v`, err)
	}
	options := syncOptions{
		Source:    sd,
		Target:    td,
		Timestamp: "timestamp",
		Noop:      false,
	}
	syncError := syncronizeDirectories(options)
	if syncError != nil {
		t.Errorf(`syncronizeDirectories error: %v`, syncError)
	}
	sync, err := directoriesAreSynchronized(options.Source, options.Target, []string{".resync"})
	if err != nil {
		t.Errorf(`unexpected error comparing dirs: %v`, err)
	}
	if !sync {
		t.Errorf(`dirs not synchronized after resync`)
	}
}
