package main

import (
	"errors"
	"log"
	"strings"

	"github.com/drone-plugins/drone-manifest/command"
)

type (
	Repo struct {
		Owner  string
		Name   string
		Branch string
	}

	Build struct {
		Path     string
		Tag      string
		Event    string
		Number   int
		Commit   string
		Ref      string
		Branch   string
		Author   string
		Pull     string
		Message  string
		DeployTo string
		Status   string
		Link     string
		Started  int64
		Created  int64
		Tags     []string
	}

	Job struct {
		Started int64
	}

	Config struct {
		Username      string
		Password      string
		Platforms     []string
		Target        string
		Template      string
		Spec          string
		IgnoreMissing bool
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Job    Job
		Config Config
	}
)

func (p *Plugin) Exec() error {
	opts := make([]command.Option, 0)

	if p.Config.Username == "" {
		return errors.New("you must provide a username")
	} else {
		opts = append(opts, command.WithUsername(p.Config.Username))
	}

	if p.Config.Password == "" {
		return errors.New("you must provide a password")
	} else {
		opts = append(opts, command.WithPassword(p.Config.Password))
	}

	if p.Config.Spec != "" {
		spec, err := RenderTrim(p.Config.Spec, p)

		if err != nil {
			return err
		}

		opts = append(opts, command.WithSpec(spec))

		log.Printf(
			"pushing by spec",
		)
	} else {
		if len(p.Config.Platforms) == 0 {
			return errors.New("you must provide platforms")
		} else {
			opts = append(opts, command.WithPlatforms(p.Config.Platforms))
		}

		if p.Config.Target == "" {
			return errors.New("you must provide a target")
		} else {
			opts = append(opts, command.WithTarget(p.Config.Target))
		}

		if p.Config.Template == "" {
			return errors.New("you must provide a template")
		} else {
			opts = append(opts, command.WithTemplate(p.Config.Template))
		}

		log.Printf(
			"pushing %s to %s for %s",
			p.Config.Template,
			p.Config.Target,
			strings.Join(p.Config.Platforms, ", "),
		)
	}

	if p.Config.IgnoreMissing {
		opts = append(opts, command.IgnoreMissing())
	}

	if p.Build.Path != "" {
		opts = append(opts, command.WithPath(p.Build.Path))
	}

	return command.New(opts...).Exec()
}
