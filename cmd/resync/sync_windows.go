// +build windows

package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/enr/go-commons/environment"
	"github.com/enr/runcmd"
	"golang.org/x/text/encoding/charmap"
)

/*
Sync using Robocopy if available, otherwise rmdir+mkdir+xcopy.
*/
func syncronizeDirectories(options syncOptions) error {
	source := options.Source
	target := options.Target
	noop := options.Noop
	robocopy := environment.Which("robocopy")
	if robocopy != "" {
		ui.Confidentialf(`Sync using robocopy %s`, robocopy)
		return syncUsingRobocopy(source, target, noop)
	}
	return syncUsingXcopy(source, target, noop)
}

/*

Sync using Robocopy tool.

Robocopy return codes
---
0	No files were copied. No failure was encountered. No files were mismatched. The files already exist in the destination directory; therefore, the copy operation was skipped.
1	All files were copied successfully.
2	There are some additional files in the destination directory that are not present in the source directory. No files were copied.
3	Some files were copied. Additional files were present. No failure was encountered.
5	Some files were copied. Some files were mismatched. No failure was encountered.
6	Additional files and mismatched files exist. No files were copied and no failures were encountered. This means that the files already exist in the destination directory.
7	Files were copied, a file mismatch was present, and additional files were present.
8	Several files did not copy.

Robocopy Logging Options
---
/L :: List only - don't copy, timestamp or delete any files.
/X :: report all eXtra files, not just those selected.
/V :: produce Verbose output, showing skipped files.
/TS :: include source file Time Stamps in the output.
/FP :: include Full Pathname of files in the output.
/BYTES :: Print sizes as bytes.
/NS :: No Size - don't log file sizes.
/NC :: No Class - don't log file classes.
/NFL :: No File List - don't log file names.
/NDL :: No Directory List - don't log directory names.
/NP :: No Progress - don't display percentage copied.
/ETA :: show Estimated Time of Arrival of copied files.
/LOG:file :: output status to LOG file (overwrite existing log).
/LOG+:file :: output status to LOG file (append to existing log).
/UNILOG:file :: output status to LOG file as UNICODE (overwrite existing log).
/UNILOG+:file :: output status to LOG file as UNICODE (append to existing log).
/TEE :: output to console window, as well as the log file.
/NJH :: No Job Header.
/NJS :: No Job Summary.
/UNICODE :: output status as UNICODE.
/R:0 :: 0 retries for read/write failures
/W:0 :: 0 seconds between retries
*/
func syncUsingRobocopy(source string, target string, noop bool) error {
	// cmd args, used to set output encoding to unicode (/U)
	cmdArgs := []string{
		// "/U",
		"/C",
	}

	robocopy := environment.Which("robocopy")
	//   /UNICODE :: output status as UNICODE.
	// otherwise exit code is evalued as "3"
	robocopyArgs := []string{
		robocopy,
		source,
		target,
		"/MIR",
		// "/UNICODE",
		"/XF",
		".resync",
	}
	if noop {
		//  /L :: List only - don't copy, timestamp or delete any files.
		ui.Warn("NOOP MODE:")
		robocopyArgs = append(robocopyArgs, "/L")
	}
	// %ComSpec%
	executable := "cmd"
	executableArgs := append(cmdArgs, robocopyArgs...)
	ui.Lifecyclef("%s %s", executable, strings.Join(executableArgs, " "))
	command := &runcmd.Command{
		Exe:  executable,
		Args: executableArgs,
	}
	res := command.Run()
	d := charmap.CodePage850.NewDecoder()
	bytes, err := d.Bytes(res.Stdout().Bytes())
	if err != nil {
		ui.Errorf("error processing output %v", err)
		return err
	}
	if noop {
		ui.Lifecycle(string(bytes))
	} else {
		ui.Confidential(string(bytes))
	}
	ui.Confidentialf("command %s. Exit code %d", command, res.ExitStatus())
	if res.ExitStatus() > 5 {
		ui.Errorf("error executing %s. Exit code %d Output:", command, res.ExitStatus())
		ui.Errorf(res.Stderr().String())
		return res.Error()
	}
	return nil
}

/*

Sync using rmdir + mkdir + xcopy.

Xcopy options
---
/E : This flag causes all folders and sub-folders to be copied
/D : This flag causes a DATE comparison to be made, only copying items that are newer than the DESTINATION item.
If the DESTINATION is older, or does not contain the file, then it will be copied.
/C : This flag tells XCOPY to continue if it encounters an error - typically errors occur with read-only files,
or files that have protected permissions
/H = Copies hidden and system files also
/R = Overwrites read-only files
/S = Copies directories and subdirectories
/Y = Overwrites existing files without asking
/O	Copies file ownership and ACL information.
*/
func syncUsingXcopy(source string, target string, noop bool) error {
	var err error
	// move or rm?
	backup := true
	if backup {
		backupDir, _ := normalizePath(target)
		err = os.MkdirAll(backupDir, os.ModePerm)
		// verificare se serve davvero o se c'e' opzione per xcopy
		// move source_dir destination
		// /y : Suppresses prompting to confirm you want to overwrite an existing destination file.
		if err == nil {
			err = execute(noop, "move", []string{target, backupDir, "/Y"})
		}
	} else {
		// /S to delete a non empty directory.
		// To delete directory in quiet mode, without being asked for confirmation, we can use /Q switch.
		// rmdir /Q /S nonemptydir
		err = execute(noop, "rmdir", []string{target, "/Q", "/S"})
	}
	if err != nil {
		ui.Errorf("error deleting directory ERROR=%v\n", err)
		return err
	}
	err = execute(noop, "mkdir", []string{target})
	if err != nil {
		ui.Errorf("error creating directory ERROR=%v\n", err)
		return err
	}
	xcopyArgs := []string{source, target, "/H", "/S", "/E", "/Y"}
	if noop {
		//  /L :: List only - don't copy, timestamp or delete any files.
		ui.Warn("NOOP MODE:")
		xcopyArgs = append(xcopyArgs, "/L")
	}
	return execute(noop, "xcopy", xcopyArgs)
}

// windows tools want dir path without the final "\".
func normalizeDirsPath(dirpath string) (string, error) {
	return normalizePath(dirpath)
}

func normalizePath(dirpath string) (string, error) {
	p, err := filepath.Abs(dirpath)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(filepath.FromSlash(p), string(filepath.Separator)), nil
}
