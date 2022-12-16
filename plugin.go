package main

import (
	"fmt"

	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/drone/drone-template-lib/template"
	"github.com/pkg/errors"
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
		Insecure      bool
		Platforms     []string
		Target        string
		Template      string
		Spec          string
		IgnoreMissing bool
		Dump          bool
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Job    Job
		Config Config
	}
)

func mainfestToolPath() string {
	if runtime.GOOS == "windows" {
		return "C:/bin/manifest-tool.exe"
	}

	return "/bin/manifest-tool"
}

func (p *Plugin) Exec() error {
	args := []string{}

	if p.Config.Username == "" && p.Config.Password != "" {
		return errors.New("you must provide a username")
	} else {
		args = append(args, fmt.Sprintf("--username=%s", p.Config.Username))
	}

	if p.Config.Password == "" && p.Config.Username != "" {
		return errors.New("you must provide a password")
	} else {
		args = append(args, fmt.Sprintf("--password=%s", p.Config.Password))
	}

	if p.Config.Insecure {
		args = append(args, "--insecure")
	}

	args = append(args, "push")

	if p.Config.Spec != "" {
		var raw []byte
		// if spec is not a valid file, assume inlining
		if _, err := os.Stat(p.Config.Spec); os.IsNotExist(err) {
			raw = []byte(p.Config.Spec)
		} else { // otherwise read it
			raw, err = os.ReadFile(p.Config.Spec)

			if err != nil {
				return errors.Wrap(err, "failed to read template")
			}
		}

		spec, err := template.RenderTrim(string(raw), p)

		if err != nil {
			return errors.Wrap(err, "failed to render template")
		}

		tmpfile, err := os.CreateTemp(p.Build.Path, "manifest-")

		if err != nil {
			return errors.Wrap(err, "failed to create tempfile")
		}

		defer os.Remove(tmpfile.Name())

		if _, err := tmpfile.Write([]byte(spec)); err != nil {
			return errors.Wrap(err, "failed to write tempfile")
		}

		if err := tmpfile.Close(); err != nil {
			return errors.Wrap(err, "failed to close tempfile")
		}

		if p.Config.Dump {
			println(spec)
		}

		args = append(args, "from-spec")
		args = append(args, tmpfile.Name())

		log.Printf(
			"pushing by spec",
		)
	} else {
		args = append(args, "from-args")

		if len(p.Config.Platforms) == 0 {
			return errors.New("you must provide platforms")
		} else {
			args = append(args, fmt.Sprintf("--platforms=%s", strings.Join(p.Config.Platforms, ",")))
		}

		if p.Config.Target == "" {
			return errors.New("you must provide a target")
		} else {
			args = append(args, fmt.Sprintf("--target=%s", p.Config.Target))
		}

		if p.Config.Template == "" {
			return errors.New("you must provide a template")
		} else {
			args = append(args, fmt.Sprintf("--template=%s", p.Config.Template))
		}

		log.Printf(
			"pushing %s to %s for %s",
			p.Config.Template,
			p.Config.Target,
			strings.Join(p.Config.Platforms, ", "),
		)
	}

	if p.Config.IgnoreMissing {
		args = append(args, "--ignore-missing")
	}

	cmd := exec.Command(
		mainfestToolPath(),
		args...,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if p.Build.Path != "" {
		cmd.Dir = p.Build.Path
	}

	return cmd.Run()
}
