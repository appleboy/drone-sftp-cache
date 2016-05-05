package sftp

import (
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
	err := c.sftp.Mkdir(filepath.Dir(p))
	if err != nil {
		return nil
	}
	dst, err := c.sftp.Create(p)
	if err != nil {
		return err
	}
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
