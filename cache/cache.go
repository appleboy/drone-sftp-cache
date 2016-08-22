package cache

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Cache implements operations for caching files.
type Cache interface {
	List(string) ([]os.FileInfo, error)
	Get(string) (io.ReadCloser, error)
	Put(string, time.Duration, io.Reader) error
	Remove(string) error
}

// Rebuild is a helper function that pushes the archived file to the cache.
func Rebuild(c Cache, src, dst string) error {
	r, w := io.Pipe()
	defer func() {
		w.Close()
		r.Close()
	}()

	c1 := make(chan error)
	c2 := make(chan error)

	go func() {
		c1 <- archive(src, w)
		w.Close()
	}()
	go func() {
		c2 <- c.Put(dst, 0, r)
		r.Close()
	}()

	err1 := <-c1
	err2 := <-c2
	if err1 != nil {
		return err1
	}
	return err2
}

// Restore is a helper function that fetches the archived file from the cache
// and restores to the host machine's file system.
func Restore(c Cache, src, dst string) error {
	rc, err := c.Get(src)
	if err != nil {
		return err
	}
	defer rc.Close()

	return extract(dst, rc)
}

//
// NOTE below are alternate implementations of the above functions that use
// the tar command for building archives.
//

// RebuildCmd is a helper function that pushes the archived file to the cache.
func RebuildCmd(c Cache, src, dst string) (err error) {

	src = filepath.Clean(src)
	src, err = filepath.Abs(src)
	if err != nil {
		return fmt.Errorf("%s, absolute path failed for %s", err, src)
	}

	// create a temporary file for the archive
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return fmt.Errorf("%s, temp dir %s failed to be created", err, os.TempDir())
	}
	tar := filepath.Join(dir, "archive.tar")

	// run archive command
	cmd := exec.Command("tar", "-cf", tar, src)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s, tar cmd failed for file %s and src %s", err, tar, src)
	}

	// upload file to server
	f, err := os.Open(tar)
	if err != nil {
		return fmt.Errorf("%s, read failed for file %s", err, tar)
	}
	defer f.Close()
	return c.Put(dst, 0, f)
}

// RestoreCmd is a helper function that fetches the archived file from the cache
// and restores to the host machine's file system.
func RestoreCmd(c Cache, src, dst string) error {
	rc, err := c.Get(src)
	if err != nil {
		return err
	}
	defer rc.Close()

	// create temp file for archive
	temp, err := ioutil.TempFile("", "")
	if err != nil {
		return err
	}
	defer func() {
		temp.Close()
		os.Remove(temp.Name())
	}()

	// download archive to temp file
	if _, err := io.Copy(temp, rc); err != nil {
		return err
	}

	// cleanup after ourself
	temp.Close()

	// run extraction command
	cmd := exec.Command("tar", "-xf", temp.Name(), "-C", "/")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
