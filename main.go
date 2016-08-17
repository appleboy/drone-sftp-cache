package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/drone-plugins/drone-sftp-cache/cache/sftp"

	"github.com/urfave/cli"
	_ "github.com/joho/godotenv/autoload"
)

var version string // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "sftp cache plugin"
	app.Usage = "use the sftp cache plugin to cache and restore build artifacts"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{

		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository full name",
			EnvVar: "DRONE_REPO",
		},
		cli.StringFlag{
			Name:   "repo.branch",
			Usage:  "repository default branch",
			EnvVar: "DRONE_REPO_BRANCH",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "repository branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},

		cli.StringSliceFlag{
			Name:   "mount",
			Usage:  "cache directories",
			EnvVar: "PLUGIN_MOUNT",
		},
		cli.BoolFlag{
			Name:   "rebuild",
			Usage:  "rebuild the cache directories",
			EnvVar: "PLUGIN_REBUILD",
		},
		cli.BoolFlag{
			Name:   "restore",
			Usage:  "restore the cache directories",
			EnvVar: "PLUGIN_RESTORE",
		},

		// private variables that should be provided via secrets.

		cli.StringFlag{
			Name:   "server",
			Usage:  "sftp server",
			EnvVar: "SFTP_CACHE_SERVER,PLUGIN_SERVER",
		},
		cli.StringFlag{
			Name:   "path",
			Usage:  "sftp server path",
			EnvVar: "SFTP_CACHE_PATH,PLUGIN_PATH",
			Value:  "/var/lib/cache/drone",
		},
		cli.StringFlag{
			Name:   "username",
			Usage:  "sftp username",
			EnvVar: "SFTP_CACHE_USERNAME,PLUGIN_USERNAME",
			Value:  "root",
		},
		cli.StringFlag{
			Name:   "password",
			Usage:  "sftp password",
			EnvVar: "SFTP_CACHE_PASSWORD,PLUGIN_PASSWORD",
		},
		cli.StringFlag{
			Name:   "key",
			Usage:  "sftp private key",
			EnvVar: "SFTP_CACHE_PRIVATE_KEY,PLUGIN_KEY",
		},
	}

	app.Run(os.Args)
}

func run(c *cli.Context) {
	plugin := Plugin{
		Mount:   c.StringSlice("mount"),
		Path:    c.String("path"),
		Repo:    c.String("repo.name"),
		Default: c.String("repo.branch"),
		Branch:  c.String("commit.branch"),
	}

	sftp, err := sftp.New(
		c.String("server"),
		c.String("username"),
		c.String("password"),
		c.String("key"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer sftp.(io.Closer).Close()

	if c.Bool("rebuild") {
		now := time.Now()
		err = plugin.Rebuild(sftp)
		log.Printf("cache built in %v", time.Since(now))
	}
	if c.Bool("restore") {
		now := time.Now()
		err = plugin.Restore(sftp)
		log.Printf("cache restored in %v", time.Since(now))
	}

	if err != nil {
		log.Println(err) // this plugins does not fail on error
	}
}
