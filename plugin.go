package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"path/filepath"

	"github.com/drone-plugins/drone-sftp-cache/cache"
)

// Plugin for caching directories to an SFTP server.
type Plugin struct {
	Mount   []string
	Path    string
	Repo    string
	Branch  string
	Default string // default master branch
}

// Rebuild the remote cache from the local environment.
func (p Plugin) Rebuild(c cache.Cache) error {
	for _, mount := range p.Mount {
		hash := hasher(mount, p.Branch)
		path := filepath.Join(p.Path, p.Repo, hash)

		log.Printf("archiving directory <%s> to remote cache <%s>", mount, path)

		err := cache.Rebuild(c, mount, path)
		if err != nil {
			return err
		}
	}
	return nil
}

// Restore the local environment from the remote cache.
func (p Plugin) Restore(c cache.Cache) error {
	for _, mount := range p.Mount {
		hash := hasher(mount, p.Branch)
		path := filepath.Join(p.Path, p.Repo, hash)

		log.Printf("restoring directory <%s> from remote cache <%s>", mount, path)

		err := cache.Restore(c, path, mount)
		if err != nil {

			// this is fallback code to restore from the projects default branch.
			// hash = hasher(mount, "master")
			// path = filepath.Join(p.Path, p.Repo, hash)
			// log.Printf("restoring directory %s from remote cache, using default branch", mount)
			// if xerr := cache.Restore(c, path, mount); xerr != nil {
			return err
			// }
		}
	}
	return nil
}

// helper function to hash a file name based on path and branch.
func hasher(mount, branch string) string {
	parts := []string{mount, branch}

	// calculate the hash using the branch
	h := md5.New()
	for _, part := range parts {
		io.WriteString(h, part)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
