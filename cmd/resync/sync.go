package main

import (
	"fmt"
	"strings"

	"github.com/enr/go-files/files"
)

type syncOptions struct {
	Source    string
	Target    string
	Timestamp string
	Noop      bool
}

// Returns a "raw" syncOption object initialized with only Source and Target fields.
func rawSyncOptions(s string, t string) (syncOptions, error) {
	options := syncOptions{}
	as := strings.TrimSpace(s)
	if as == "" {
		return options, &appError{code: 3, message: "Error empty source", cause: nil}
	}
	source, err := normalizeDirsPath(as)
	if err != nil {
		return options, &appError{code: 3, message: fmt.Sprintf("Error reading source path %s: %v", as, err), cause: err}
	}
	at := strings.TrimSpace(t)
	if at == "" {
		return options, &appError{code: 3, message: "Error empty target", cause: nil}
	}
	target, err := normalizeDirsPath(at)
	if err != nil {
		return options, &appError{code: 3, message: fmt.Sprintf("Error reading target path %s: %v", at, err), cause: err}
	}
	ui.Confidentialf("Synchronize source=%s target=%s", source, target)
	if !files.IsDir(source) {
		return options, &appError{code: 4, message: fmt.Sprintf("Source should be a directory %s", source), cause: nil}
	}
	if !files.IsDir(target) {
		return options, &appError{code: 4, message: fmt.Sprintf("Target should be a directory %s", target), cause: nil}
	}
	if isInside(source, target) || isInside(target, source) {
		return options, &appError{code: 6, message: fmt.Sprintf("Overlapping directories: %s %s", source, target), cause: nil}
	}
	options.Source = source
	options.Target = target
	return options, nil
}

func writeReport(options syncOptions, syncError error) error {
	source := options.Source
	target := options.Target
	timestamp := options.Timestamp
	noop := options.Noop
	if syncError != nil {
		ui.Errorf("Error synchronizing: %v", syncError)
	}
	ui.Confidentialf("Sync timestamp: %s\n", timestamp)
	var result string
	if noop {
		result = "[skip]"
	} else {
		result = "[FAIL]"
		if syncError == nil {
			result = "[ok  ]"
		}
	}
	line := fmt.Sprintf("%s %s %s %s\n", timestamp, result, source, target)
	ui.Confidentialf("Writing .resync line:\n%s", line)
	var err error
	if !noop {
		err = writeReportFile(source, line)
		if err != nil {
			ui.Errorf("Error write report %s: %v", source, err)
			return err
		}
		err = writeReportFile(target, line)
		if err != nil {
			ui.Errorf("Error write report %s: %v", source, err)
			return err
		}
	}
	return nil
}
