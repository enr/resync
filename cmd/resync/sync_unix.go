// +build darwin freebsd linux netbsd openbsd

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/enr/go-commons/environment"
)

func syncronizeDirectories(options syncOptions) error {
	source := options.Source
	target := options.Target
	noop := options.Noop
	ui.Confidentialf(`Synchronize(noop=%t) "%s" -> "%s"`, noop, source, target)
	// rsync --modify-window=1 -rltDv --delete --copy-links /home/enrico/Projects/ /media/HP\ v135w/Projects/
	rsync := environment.Which("rsync")
	if rsync == "" {
		ui.Errorf("rsync not found in path. exit\n")
		os.Exit(1)
	}
	ui.Confidentialf("rsync found %s\n", rsync)
	rsyncShortOpts := "-rltDv"
	if noop {
		// -n, --dry-run: perform a trial run with no changes made
		rsyncShortOpts = rsyncShortOpts + "n"
		ui.Warn("NOOP MODE:")
	}
	rsyncArgs := []string{
		"--modify-window=1",
		rsyncShortOpts,
		"--delete",
		"--copy-links",
		"--exclude",
		".resync",
		source,
		target,
	}
	return execute(noop, rsync, rsyncArgs)
}

// ensure that directory path ends with "/"
func normalizeDirsPath(dirpath string) (string, error) {
	p, err := normalizePath(dirpath)
	if err != nil {
		return "", err
	}
	if strings.HasSuffix(p, "/") {
		return p, nil
	}
	return fmt.Sprintf("%s/", p), nil
}

// normalize for unix: absolute path ending with "/" if directory
func normalizePath(dirpath string) (string, error) {
	p, err := filepath.Abs(dirpath)
	if err != nil {
		return "", err
	}
	p = filepath.ToSlash(p)
	return p, nil
}
