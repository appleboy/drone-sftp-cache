package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"regexp"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/appleboy/drone-sftp-cache/cache"
	"github.com/appleboy/drone-sftp-cache/cache/sftp"
)

var skipRe = regexp.MustCompile(`\[(?i:cache *skip|skip *cache)\]`)

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
	Default      string
	Message      string
}

// Exec executes the plugin.
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
		// skip the restore if any case-insensitive combination of the words "skip" and "cache"
		skipMatch := skipRe.FindString(p.Message)
		if len(skipMatch) > 0 {
			logrus.Printf("skip restore cache. %s found in '%s'", skipMatch, p.Message)
			return nil
		}

		now := time.Now()
		err = p.ProcessRestore(sftp)
		logrus.Printf("cache restored in %v", time.Since(now))
	}

	if err != nil {
		logrus.Println(err)
	}

	return nil
}

// ProcessRebuild rebuild the remote cache from the local environment.
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

// ProcessRestore restore the local environment from the remote cache.
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
			return err
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
