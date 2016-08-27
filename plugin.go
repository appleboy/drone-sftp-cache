package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/drone-plugins/drone-sftp-cache/cache"
	"github.com/drone-plugins/drone-sftp-cache/cache/sftp"
)

// Plugin for caching directories to an SFTP server.
type Plugin struct {
	IgnoreBranch bool
	Rebuild      bool
	Restore      bool
	Server       string
	Username     string
	Password     string
	Key          string
	Mount        []string
	Path         string
	Repo         string
	Branch       string
	Default      string // default master branch
}

func (p *Plugin) Exec() error {
	sftp, err := sftp.New(
		p.Server,
		p.Username,
		p.Password,
		p.Key,
	)

	if err != nil {
		return err
	}

	defer sftp.(io.Closer).Close()

	if p.Rebuild {
		now := time.Now()
		err = p.ProcessRebuild(sftp)
		logrus.Printf("cache built in %v", time.Since(now))
	}

	if p.Restore {
		now := time.Now()
		err = p.ProcessRestore(sftp)
		logrus.Printf("cache restored in %v", time.Since(now))
	}

	if err != nil {
		logrus.Println(err)
	}

	return nil
}

// Rebuild the remote cache from the local environment.
func (p Plugin) ProcessRebuild(c cache.Cache) error {
	for _, mount := range p.Mount {
		var hash string
		if p.IgnoreBranch {
			hash = hasher(mount)
		} else {
			hash = hasher(mount, p.Branch)
		}
		path := filepath.Join(p.Path, p.Repo, hash)

		log.Printf("archiving directory <%s> to remote cache <%s>", mount, path)

		err := cache.RebuildCmd(c, mount, path)
		if err != nil {
			return err
		}
	}
	return nil
}

// Restore the local environment from the remote cache.
func (p Plugin) ProcessRestore(c cache.Cache) error {
	for _, mount := range p.Mount {
		var hash string
		if p.IgnoreBranch {
			hash = hasher(mount)
		} else {
			hash = hasher(mount, p.Branch)
		}
		path := filepath.Join(p.Path, p.Repo, hash)

		log.Printf("restoring directory <%s> from remote cache <%s>", mount, path)

		err := cache.RestoreCmd(c, path, mount)
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
func hasher(args ...string) string {
	// calculate the hash using the branch
	h := md5.New()
	for _, part := range args {
		io.WriteString(h, part)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
