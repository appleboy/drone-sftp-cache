package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/crypto/ssh"

	"github.com/pkg/sftp"
)

type Plugin struct {
	Mount    []string
	Server   string
	Path     string
	Username string
	Password string
	Key      string
	Repo     string
	Branch   string
	Rebuild  bool
	Restore  bool
}

func (p Plugin) Exec() error {

	config := &ssh.ClientConfig{
		User: p.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(p.Password),
		},
	}

	// private key authentication takes precedence
	if p.Key != "" {
		signer, err := ssh.ParsePrivateKey([]byte(p.Key))
		if err != nil {
			return err
		}
		config.Auth[0] = ssh.PublicKeys(signer)
	}

	// create the ssh connection and client
	client, err := ssh.Dial("tcp", p.Server, config)
	if err != nil {
		return err
	}
	defer client.Close()

	// open the sftp session using the ssh connection
	sftp, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer sftp.Close()

	// if restoring, download and unpack
	for _, mount := range p.Mount {
		hash := signature(mount, p.Branch)
		path := filepath.Join(p.Path, p.Repo, hash)

		// download the remote file to the specified target path
		src, err := download(sftp, path)
		if err != nil {
			return err
		}
		defer src.Close()

		// create the temporary file where we can stream the download.
		dst, err := ioutil.TempFile("", "")
		if err != nil {
			return err
		}
		defer dst.Close()

		// copies the download stream to the temp file.
		_, err = io.Copy(dst, src)
		if err != nil {
			return err
		}

		// extract the archive to the destination directory.
		if err = restore(dst.Name()); err != nil {
			return err
		}
	}

	// walk the directory on the remote server
	// w := sftp.Walk(p.Path)
	// for w.Step() {
	// 	if w.Err() != nil {
	// 		continue
	// 	}
	// 	log.Println(w.Path())
	// }

	return nil
}

// signature is a helper function for creating a signature for cache archives.
func signature(mount, branch string) string {
	parts := []string{mount, branch}

	// calculate the hash using the branch
	h := md5.New()
	for _, part := range parts {
		io.WriteString(h, part)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func download(client *sftp.Client, path string) (io.ReadCloser, error) {
	_, err := client.Stat(path)
	if err != nil {
		return nil, err
	}
	return client.Open(path)
}

func restore(tar string) error {
	opt := untarOpts("")
	cmd := exec.Command("tar", opt, tar, "-C", "/")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func rebuild(dir, tar string) (err error) {
	opt := tarOpts("")
	cmd := exec.Command("tar", opt, tar, dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func tarOpts(archive string) string {
	switch archive {
	case "bzip", "bzip2":
		return "-cjf"
	case "gzip":
		return "-czf"
	default:
		return "-cf"
	}
}

func untarOpts(archive string) string {
	switch archive {
	case "bzip", "bzip2":
		return "-xjf"
	case "gzip":
		return "-xzf"
	default:
		return "-xf"
	}
}
