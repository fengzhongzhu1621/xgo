package shell

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"
)

var pidfile string

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Errorf("get current path failed. Error:%s", err.Error())
	}
	pidfile = cwd + "/pid/" + filepath.Base(os.Args[0]) + ".pid"
}

// SavePid TODO
func SavePid() error {
	if err := WritePid(); err != nil {
		return fmt.Errorf("write pid file failed. err:%s", err.Error())
	}

	return nil
}

// SetPidfilePath sets the pidfile path.
func SetPidfilePath(p string) {
	pidfile = p
}

// WritePid the pidfile based on the flag. It is an error if the pidfile hasn't
// been configured.
func WritePid() error {
	if pidfile == "" {
		return fmt.Errorf("pidfile is not set")
	}

	if err := os.MkdirAll(filepath.Dir(pidfile), os.FileMode(0755)); err != nil {
		return err
	}

	file, err := AtomicFileNew(pidfile, os.FileMode(0644))
	if err != nil {
		return fmt.Errorf("error opening pidfile %s: %s", pidfile, err)
	}
	defer file.Close() // in case we fail before the explicit close

	_, err = fmt.Fprintf(file, "%d", os.Getpid())
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

// ReadPid the pid from the configured file. It is an error if the pidfile hasn't
// been configured.
func ReadPid() (int, error) {
	if pidfile == "" {
		return 0, fmt.Errorf("pidfile is empty")
	}

	d, err := ioutil.ReadFile(pidfile)
	if err != nil {
		return 0, err
	}

	pid, err := strconv.Atoi(string(bytes.TrimSpace(d)))
	if err != nil {
		return 0, fmt.Errorf("error parsing pid from %s: %s", pidfile, err)
	}

	return pid, nil
}
