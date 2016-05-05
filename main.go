package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
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
			Name:   "repo",
			Usage:  "repository full name",
			EnvVar: "DRONE_REPO",
		},
		cli.StringFlag{
			Name:   "branch",
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
			EnvVar: "SFTP_CACHE_SERVER",
		},
		cli.StringFlag{
			Name:   "path",
			Usage:  "sftp server path",
			EnvVar: "SFTP_CACHE_PATH",
			Value:  "/var/lib/cache/drone",
		},
		cli.StringFlag{
			Name:   "username",
			Usage:  "sftp username",
			EnvVar: "SFTP_CACHE_USERNAME",
			Value:  "root",
		},
		cli.StringFlag{
			Name:   "password",
			Usage:  "sftp password",
			EnvVar: "SFTP_CACHE_PASSWORD",
		},
		cli.StringFlag{
			Name:   "key",
			Usage:  "sftp private key",
			EnvVar: "SFTP_CACHE_PRIVATE_KEY",
		},
	}

	app.Run(os.Args)
}

func run(c *cli.Context) {
	plugin := Plugin{
		Mount:    c.StringSlice("mount"),
		Path:     c.String("path"),
		Server:   c.String("server"),
		Username: c.String("username"),
		Password: c.String("password"),
		Key:      c.String("key"),
		Repo:     c.String("repo"),
		Branch:   c.String("branch"),
		Rebuild:  c.Bool("rebuild"),
		Restore:  c.Bool("restore"),
	}

	if err := plugin.Exec(); err != nil {
		log.Println(err) // this plugins does not fail on error
	}
}
