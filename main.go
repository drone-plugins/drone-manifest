package main

import (
	"fmt"
	"log"
	"os"

	"github.com/drone-plugins/drone-manifest/tagging"
	"github.com/urfave/cli"
)

var (
	version = "unknown"
)

func main() {
	app := cli.NewApp()
	app.Name = "manifest plugin"
	app.Usage = "manifest plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "username",
			Usage:  "username for registry",
			EnvVar: "PLUGIN_USERNAME,MANIFEST_USERNAME,DOCKER_USERNAME",
		},
		cli.StringFlag{
			Name:   "password",
			Usage:  "password for registry",
			EnvVar: "PLUGIN_PASSWORD,MANIFEST_PASSWORD,DOCKER_PASSWORD",
		},
		cli.StringSliceFlag{
			Name:   "platforms",
			Usage:  "platforms for manifests",
			EnvVar: "PLUGIN_PLATFORMS",
		},
		cli.StringFlag{
			Name:   "target",
			Usage:  "target for manifests",
			EnvVar: "PLUGIN_TARGET",
		},
		cli.StringFlag{
			Name:   "template",
			Usage:  "template for manifests",
			EnvVar: "PLUGIN_TEMPLATE",
		},
		cli.StringFlag{
			Name:   "spec",
			Usage:  "path to manifest spec",
			EnvVar: "PLUGIN_SPEC",
		},
		cli.BoolFlag{
			Name:   "ignore-missing",
			Usage:  "ignore missing images",
			EnvVar: "PLUGIN_IGNORE_MISSING",
		},
		cli.StringSliceFlag{
			Name:     "tags",
			Usage:    "list of additional tags",
			Value:    &cli.StringSlice{},
			EnvVar:   "PLUGIN_TAG,PLUGIN_TAGS",
			FilePath: ".tags",
		},
		cli.BoolFlag{
			Name:   "tags.auto",
			Usage:  "automatically build tags",
			EnvVar: "PLUGIN_DEFAULT_TAGS,PLUGIN_AUTO_TAG",
		},
		cli.StringFlag{
			Name:   "path",
			Usage:  "git clone path",
			EnvVar: "DRONE_WORKSPACE",
		},
		cli.StringFlag{
			Name:   "repo.owner",
			Usage:  "repository owner",
			EnvVar: "DRONE_REPO_OWNER",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "repo.branch",
			Usage:  "repository default branch",
			EnvVar: "DRONE_REPO_BRANCH",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA",
			Value:  "00000000",
		},
		cli.StringFlag{
			Name:   "commit.ref",
			Usage:  "git commit ref",
			EnvVar: "DRONE_COMMIT_REF",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "git commit branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},
		cli.StringFlag{
			Name:   "commit.pull",
			Usage:  "git pull request",
			EnvVar: "DRONE_PULL_REQUEST",
		},
		cli.StringFlag{
			Name:   "build.event",
			Value:  "push",
			Usage:  "build event",
			EnvVar: "DRONE_BUILD_EVENT",
		},
		cli.IntFlag{
			Name:   "build.number",
			Usage:  "build number",
			EnvVar: "DRONE_BUILD_NUMBER",
		},
		cli.StringFlag{
			Name:   "build.status",
			Usage:  "build status",
			Value:  "success",
			EnvVar: "DRONE_BUILD_STATUS",
		},
		cli.StringFlag{
			Name:   "build.link",
			Usage:  "build link",
			EnvVar: "DRONE_BUILD_LINK",
		},
		cli.Int64Flag{
			Name:   "build.started",
			Usage:  "build started",
			EnvVar: "DRONE_BUILD_STARTED",
		},
		cli.Int64Flag{
			Name:   "build.created",
			Usage:  "build created",
			EnvVar: "DRONE_BUILD_CREATED",
		},
		cli.StringFlag{
			Name:   "build.tag",
			Usage:  "build tag",
			EnvVar: "DRONE_TAG",
		},
		cli.Int64Flag{
			Name:   "job.started",
			Usage:  "job started",
			EnvVar: "DRONE_JOB_STARTED",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Repo: Repo{
			Owner:  c.String("repo.owner"),
			Name:   c.String("repo.name"),
			Branch: c.String("repo.branch"),
		},
		Build: Build{
			Path:    c.String("path"),
			Tag:     c.String("build.tag"),
			Number:  c.Int("build.number"),
			Event:   c.String("build.event"),
			Status:  c.String("build.status"),
			Commit:  c.String("commit.sha"),
			Ref:     c.String("commit.ref"),
			Branch:  c.String("commit.branch"),
			Pull:    c.String("commit.pull"),
			Started: c.Int64("build.started"),
			Created: c.Int64("build.created"),
			Tags:    c.StringSlice("tags"),
		},
		Job: Job{
			Started: c.Int64("job.started"),
		},
		Config: Config{
			Username:      c.String("username"),
			Password:      c.String("password"),
			Platforms:     c.StringSlice("platforms"),
			Target:        c.String("target"),
			Template:      c.String("template"),
			Spec:          c.String("spec"),
			IgnoreMissing: c.Bool("ignore-missing"),
		},
	}

	if c.Bool("tags.auto") {
		if tagging.UseDefaultTag(c.String("commit.ref"), c.String("repo.branch")) {
			plugin.Build.Tags = tagging.DefaultTags(c.String("commit.ref"))
		} else {
			log.Printf("skipping automated tags for %s", c.String("commit.ref"))
			return nil
		}
	}

	return plugin.Exec()
}
