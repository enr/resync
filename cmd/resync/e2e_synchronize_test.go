package main

import (
	"fmt"
	"testing"

	"github.com/enr/runcmd"
)

func TestSynch(t *testing.T) {
	src := testdataPath
	target, err := tempDirectory()
	if err != nil {
		t.Errorf(`unexpected error creating target dir: %v`, err)
	}
	writeFileInDir(t, target, ".hidden")
	writeFileInDir(t, target, "test.txt")
	writeFileInDir(t, target, "01.txt")
	writeFileInDir(t, target, "uno/wow.txt")

	// verify directories are differents
	sync, _ := directoriesAreSynchronized(src, target, []string{".resync"})
	if sync {
		t.Errorf(`dirs already sync before resync`)
	}

	command := &runcmd.Command{
		Exe:  exePath("resync"),
		Args: []string{"-c", fmt.Sprintf("%s/resyncrc.yaml", testdataPath), src, target},
	}
	res := command.Run()
	if !res.Success() {
		fmt.Println(res.Stderr().String())
		t.Fatalf("%s: unexpected fail status=%d", command, res.ExitStatus())
	}
	expectedCode := 0
	actualCode := res.ExitStatus()
	if actualCode != expectedCode {
		t.Fatalf("%s: expected exit code %d but got %d", command, expectedCode, actualCode)
	}

	sync, err = directoriesAreSynchronized(src, target, []string{".resync"})
	if err != nil {
		t.Errorf(`unexpected error comparing dirs: %v`, err)
	}
	if !sync {
		t.Errorf(`dirs not synchronized after resync`)
	}
}
