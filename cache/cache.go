package cache

import (
	"io"
	"os"
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
