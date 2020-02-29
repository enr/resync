package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/enr/go-files/files"
)

type rawOptionsTestCase struct {
	inputSource   string
	inputTarget   string
	success       bool
	expectedError *appError
}

var rawOptionsFixture = []rawOptionsTestCase{
	{
		inputSource:   "  ",
		inputTarget:   "",
		success:       false,
		expectedError: &appError{code: 3},
	},
	{
		inputSource:   "/tmp/a/b/c/d",
		inputTarget:   "/tmp/a",
		success:       false,
		expectedError: &appError{code: 6},
	},
	{
		inputSource:   "/tmp/a/b",
		inputTarget:   "/tmp/a/b/c",
		success:       false,
		expectedError: &appError{code: 6},
	},
}

func TestRawSyncOptions(t *testing.T) {
	os.MkdirAll("/tmp/a/b/c/d", 0777)
	for _, testcase := range rawOptionsFixture {
		_, err := rawSyncOptions(testcase.inputSource, testcase.inputTarget)

		if testcase.success {
			if err != nil {
				t.Errorf(`%s -> %s unexpected error %v`, testcase.inputSource, testcase.inputTarget, err)
			}
			continue
		}
		if err == nil {
			t.Errorf(`%s -> %s expected error got nil`, testcase.inputSource, testcase.inputTarget)
		}
		if testcase.expectedError != nil {

			if ae, ok := err.(*appError); ok {
				if ae.code != testcase.expectedError.code {
					t.Errorf(`%s -> %s expected error code %d got %d`, testcase.inputSource, testcase.inputTarget, testcase.expectedError.code, ae.code)
				}
			}
		}
	}
}

func TestWriteReport(t *testing.T) {
	sd, err := tempDirectory()
	if err != nil {
		t.Errorf(`unexpected error creating temp dir: %v`, err)
	}
	td, err := tempDirectory()
	if err != nil {
		t.Errorf(`unexpected error creating temp dir: %v`, err)
	}
	options := syncOptions{
		Source:    sd,
		Target:    td,
		Timestamp: "TestWriteReport-timestamp",
		Noop:      false,
	}
	syncError := writeReport(options, nil)
	if syncError != nil {
		t.Errorf(`writeReport returns unexpected error: %v`, syncError)
	}
	resyncFilePath := path.Join(td, ".resync")
	b, err := ioutil.ReadFile(resyncFilePath)
	if err != nil {
		t.Errorf(`unexpected error reading .resync file: %v`, err)
	}
	resyncContents := string(b)
	assertStringContains(t, resyncContents, "TestWriteReport-timestamp [ok  ]")
	assertStringContains(t, resyncContents, options.Source)
	assertStringContains(t, resyncContents, options.Target)
}

func TestWriteReportNoop(t *testing.T) {
	sd, err := tempDirectory()
	if err != nil {
		t.Errorf(`unexpected error creating temp dir: %v`, err)
	}
	td, err := tempDirectory()
	if err != nil {
		t.Errorf(`unexpected error creating temp dir: %v`, err)
	}
	options := syncOptions{
		Source:    sd,
		Target:    td,
		Timestamp: "TestWriteReport-timestamp",
		Noop:      true,
	}
	syncError := writeReport(options, nil)
	if syncError != nil {
		t.Errorf(`writeReport returns unexpected error: %v`, syncError)
	}
	resyncFilePath := path.Join(td, ".resync")
	if files.Exists(resyncFilePath) {
		t.Errorf(`writeReport write .resync file in noop mode`)
	}
}
