package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"

	"github.com/hpcloud/tail"

	"github.com/enr/go-files/files"
	"github.com/enr/runcmd"
)

func timestamp() string {
	const layout = "2006.01.02 15:04:05"
	t := time.Now()
	ts := t.Local().Format(layout)
	return ts
}

func logPath() string {
	const layout = "20060102150405"
	t := time.Now()
	ts := t.Local().Format(layout)
	return fmt.Sprintf("resync-%s.log", ts)
}

func openFile(resyncFile string) (file *os.File, err error) {
	if files.Exists(resyncFile) {
		return os.OpenFile(resyncFile, os.O_APPEND|os.O_WRONLY, 0600)
	}
	return os.Create(resyncFile)
}

func writeReportFile(baseDir string, line string) error {
	resyncFile, _ := normalizePath(fmt.Sprintf("%s/.resync", baseDir))
	ui.Confidentialf("Dot resync: %s", resyncFile)
	f, err := openFile(resyncFile)
	defer f.Close()
	if err != nil {
		ui.Errorf("1. Error write report %s: %v", resyncFile, err)
		return err
	}
	if _, err = f.WriteString(line); err != nil {
		ui.Errorf("2. Error write report %s: %v", resyncFile, err)
		return err
	}
	return nil
}

func execute(noop bool, executable string, executableArgs []string) error {
	ui.Lifecyclef("%s %s", executable, strings.Join(executableArgs, " "))
	command := &runcmd.Command{
		Exe:     executable,
		Args:    executableArgs,
		Logfile: logPath(),
	}
	logFile := command.GetLogfile()
	fmt.Println("Log file " + logFile)
	var t *tail.Tail
	go func() {
		t, _ = tail.TailFile(logFile, tail.Config{Follow: true})
		for line := range t.Lines {
			fmt.Println(line.Text)
		}
	}()
	err := command.Start()

	if err != nil {
		ui.Errorf("error starting %s", command)
		return err
	}
	runningProcess = command.Process
	// exit := make(chan error, 2)
	// go func() {
	state, err := runningProcess.Wait()
	fmt.Printf("%v \n", state)
	ui.Warnf("call wait %v %v", err, state)
	if err != nil {
		ui.Errorf("error process %s", command)
		return err
	}
	err = t.Stop()
	if err != nil {
		ui.Errorf("error tail stop %s", command)
		return err
	}
	// 	exit <- err
	// }()

	// res := command.Run()
	// if noop {
	// 	ui.Lifecycle(res.Stdout().String())
	// } else {
	// 	ui.Confidential(res.Stdout().String())
	// }
	// if !res.Success() {
	// 	ui.Errorf("error executing %s. Exit code %d Output:", command, res.ExitStatus())
	// 	ui.Errorf(res.Stderr().String())
	// 	return res.Error()
	// }
	return nil
}

// In configuration file source paths are allowed to start with "~" as a shortcut for user's home.
func sourcePath(ppath string) (string, error) {
	retpath := ppath
	if strings.HasPrefix(ppath, "~") {
		dir, err := homedir.Dir()
		if err != nil {
			return "", err
		}
		relpath := strings.TrimPrefix(ppath, "~")
		retpath = filepath.FromSlash(path.Join(dir, relpath))
	}
	return retpath, nil
}

// returns true if p1 is inside p2
// expects two normalized, absolute paths
func isInside(pathUnderTest string, directory string) bool {
	p1 := strings.TrimSpace(filepath.FromSlash(pathUnderTest))
	fmt.Println("p1 " + p1)
	if p1 != "" && !strings.HasSuffix(p1, string(os.PathSeparator)) {
		p1 = p1 + string(os.PathSeparator)
	}
	fmt.Println("p1 " + p1)
	p2 := strings.TrimSpace(filepath.FromSlash(directory))
	fmt.Println("p2 " + p2)
	if p2 != "" && !strings.HasSuffix(p2, string(os.PathSeparator)) {
		p2 = p2 + string(os.PathSeparator)
	}
	fmt.Println("p2 " + p2)

	if p1 == p2 {
		return true
	}
	if p1 == "" || p2 == "" {
		return false
	}
	if strings.HasPrefix(p1, p2) {
		return true
	}
	return false
}
