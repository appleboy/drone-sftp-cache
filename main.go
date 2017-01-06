package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli"
)

// Version set at compile-time
var Version string

func main() {
	app := cli.NewApp()
	app.Name = "sftp cache plugin"
	app.Usage = "sftp cache plugin"
	app.Action = run
	app.Version = Version
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
		cli.BoolFlag{
			Name:   "ignore_branch",
			Usage:  "ignore branch name on hash value",
			EnvVar: "PLUGIN_IGNORE_BRANCH",
		},
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
		cli.StringFlag{
			Name:  "env-file",
			Usage: "source env file",
		},
		cli.StringFlag{
			Name:   "commit.message",
			Usage:  "commit message",
			EnvVar: "DRONE_COMMIT_MESSAGE",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Plugin{
		IgnoreBranch: c.Bool("ignore_branch"),
		Rebuild:      c.Bool("rebuild"),
		Restore:      c.Bool("restore"),
		Server:       c.String("server"),
		Username:     c.String("username"),
		Password:     c.String("password"),
		Key:          c.String("key"),
		Mount:        c.StringSlice("mount"),
		Path:         c.String("path"),
		Repo:         c.String("repo.name"),
		Default:      c.String("repo.branch"),
		Branch:       c.String("commit.branch"),
		Message:      c.String("commit.message"),
	}

	return plugin.Exec()
}
