package main

import (
	"fmt"
)

const (
	// error messages
	errorMessageArgsNumber string = "args error. expecting at least a source and a target"

	// exit code
	genericErrorExitCode          int = 1
	argsNumberErrorExitCode       int = 4
	configurationFileNotFoundCode int = 3
)

type appError struct {
	code    int
	message string
	cause   error
}

func (e *appError) Error() string {
	return fmt.Sprintf("%d - %s", e.code, e.message)
}

func manageOptionsError(err error) (bool, int) {
	if err == nil {
		return true, 0
	}
	exitCode := genericErrorExitCode
	switch e := err.(type) {
	case *appError:
		ui.Errorf(e.message)
		if e.cause != nil {
			ui.Errorf("Cause: %v", e.cause)
		}
		exitCode = e.code
	default:
		ui.Errorf("%v", e)
	}
	return false, exitCode
}
