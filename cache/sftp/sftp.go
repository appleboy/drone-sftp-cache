package sftp

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/drone-plugins/drone-sftp-cache/cache"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// cacher is an SFTP implementation of the Cache.
type cacher struct {
	sftp *sftp.Client
	ssh  *ssh.Client
}

// List returns a list of all files at the defined path.
func (c *cacher) List(root string) ([]os.FileInfo, error) {
	var files []os.FileInfo

	f := c.sftp.Walk(root)
	for f.Step() {
		if f.Err() != nil {
			continue
		}
		files = append(files, f.Stat())
	}
	return files, nil
}

// Get returns an io.Reader for reading the contents of the file.
func (c *cacher) Get(p string) (io.ReadCloser, error) {
	_, err := c.sftp.Stat(p)
	if err != nil {
		return nil, err
	}
	return c.sftp.Open(p)
}

// Put uploads the contents of the io.Reader to the SFTP server.
func (c *cacher) Put(p string, t time.Duration, src io.Reader) error {
	dir := filepath.Dir(p)

	if _, serr := c.sftp.Stat(dir); serr != nil {
		err := c.mkdirAll(dir, 0700)
		if err != nil {
			return fmt.Errorf("%s, mkdir failed for %s", err, dir)
		}
	}

	dst, err := c.sftp.Create(p)
	if err != nil {
		return fmt.Errorf("%s, sftp create failed for %s", err, p)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

// Remove removes the file from the remote SFTP server.
func (c *cacher) Remove(p string) error {
	_, err := c.sftp.Stat(p)
	if err != nil {
		return err
	}
	return c.sftp.Remove(p)
}

// Close closes the SFTP connection.
func (c *cacher) Close() error {
	if c.ssh != nil {
		c.ssh.Close()
	}
	if c.ssh != nil {
		c.sftp.Close()
	}
	return nil
}

func (c *cacher) mkdir(name string, perm os.FileMode) error {
	err := c.sftp.Mkdir(name)
	if err != nil {
		return err
	}
	return c.chmod(name, perm)
}

func (c *cacher) mkdirAll(path string, perm os.FileMode) error {
	// Fast path: if we can tell whether path is a directory or file, stop with success or error.
	dir, err := c.sftp.Stat(path)
	if err == nil {
		if dir.IsDir() {
			return nil
		}
		return err
	}

	// Slow path: make sure parent exists and then call Mkdir for path.
	i := len(path)
	for i > 0 && os.IsPathSeparator(path[i-1]) { // Skip trailing path separator.
		i--
	}

	j := i
	for j > 0 && !os.IsPathSeparator(path[j-1]) { // Scan backward over element.
		j--
	}

	if j > 1 {
		// Create parent
		err = c.mkdirAll(path[0:j-1], perm)
		if err != nil {
			return err
		}
	}

	// Parent now exists; invoke Mkdir and use its result.
	err = c.mkdir(path, perm)
	if err != nil {
		// Handle arguments like "foo/." by
		// double-checking that directory doesn't exist.
		dir, err1 := c.lstat(path)
		if err1 == nil && dir.IsDir() {
			return nil
		}
		return err
	}
	return nil
}

func (c *cacher) stat(name string) (os.FileInfo, error) {
	return c.sftp.Stat(name)
}

func (s *cacher) lstat(p string) (os.FileInfo, error) {
	return s.sftp.Lstat(p)
}

func (s *cacher) chmod(name string, mode os.FileMode) error {
	return s.sftp.Chmod(name, mode)
}

// New returns a new SFTP remote Cache implementated.
func New(server, username, password, key string) (cache.Cache, error) {
	config := &ssh.ClientConfig{
		Timeout: time.Minute * 5,
		User:    username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}

	// private key authentication takes precedence
	if key != "" {
		signer, err := ssh.ParsePrivateKey([]byte(key))
		if err != nil {
			return nil, err
		}
		config.Auth[0] = ssh.PublicKeys(signer)
	}

	// create the ssh connection and client
	client, err := ssh.Dial("tcp", server, config)
	if err != nil {
		return nil, err
	}

	// open the sftp session using the ssh connection
	sftp, err := sftp.NewClient(client)
	if err != nil {
		client.Close()
		return nil, err
	}

	return &cacher{sftp, client}, nil
}
