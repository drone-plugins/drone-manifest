// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/drone-plugins/drone-manifest/tagging"
	"github.com/drone/drone-template-lib/template"
	"github.com/urfave/cli/v2"
)

// Settings for the plugin.
type Settings struct {
	Username      string
	Password      string
	Insecure      bool
	Platforms     cli.StringSlice
	Target        string
	Template      string
	Spec          string
	IgnoreMissing bool
	Tags          cli.StringSlice
	AutoTag       bool
}

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	if p.settings.Username == "" && p.settings.Password != "" {
		return errors.New("you must provide a username")
	}

	if p.settings.Password == "" && p.settings.Username != "" {
		return errors.New("you must provide a password")
	}

	if p.settings.Spec == "" {
		if len(p.settings.Platforms.Value()) == 0 {
			return errors.New("you must provide platforms")
		}

		if p.settings.Target == "" {
			return errors.New("you must provide a target")
		}

		if p.settings.Template == "" {
			return errors.New("you must provide a template")
		}
	}

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {
	// Anonymous struct for the templating engine
	var t struct {
		Build struct {
			Tag  string
			Tags []string
		}
	}

	// Determine the tags
	t.Build.Tag = p.pipeline.Build.Tag
	if p.settings.AutoTag {
		if tagging.UseDefaultTag(p.pipeline.Commit.Ref, p.pipeline.Repo.Branch) {
			t.Build.Tags = tagging.DefaultTags(p.pipeline.Commit.Ref)
		} else {
			log.Printf("skipping automated tags for %s", p.pipeline.Commit.Ref)
			return nil
		}
	} else {
		t.Build.Tags = p.settings.Tags.Value()
	}

	args := []string{
		fmt.Sprintf("--username=%s", p.settings.Username),
		fmt.Sprintf("--password=%s", p.settings.Password),
	}

	if p.settings.Insecure {
		args = append(args, "--insecure")
	}

	args = append(args, "push")

	if p.settings.Spec != "" {
		var raw []byte
		// if spec is not a valid file, assume inlining
		if _, err := os.Stat(p.settings.Spec); os.IsNotExist(err) {
			raw = []byte(p.settings.Spec)
		} else { // otherwise read it
			raw, err = ioutil.ReadFile(p.settings.Spec)

			if err != nil {
				return fmt.Errorf("failed to read template: %w", err)
			}
		}

		spec, err := template.RenderTrim(string(raw), t)

		if err != nil {
			return fmt.Errorf("failed to render template: %w", err)
		}

		tmpfile, err := ioutil.TempFile("", "manifest-")

		if err != nil {
			return fmt.Errorf("failed to create tempfile: %w", err)
		}

		defer os.Remove(tmpfile.Name())

		if _, err := tmpfile.Write([]byte(spec)); err != nil {
			return fmt.Errorf("failed to write tempfile: %w", err)
		}

		if err := tmpfile.Close(); err != nil {
			return fmt.Errorf("failed to close tempfile: %w", err)
		}

		args = append(args, "from-spec")
		args = append(args, tmpfile.Name())
	} else {
		args = append(
			args,
			"from-args",
			fmt.Sprintf("--platforms=%s", strings.Join(p.settings.Platforms.Value(), ",")),
			fmt.Sprintf("--target=%s", p.settings.Target),
			fmt.Sprintf("--template=%s", p.settings.Template),
		)
	}

	if p.settings.IgnoreMissing {
		args = append(args, "--ignore-missing")
	}

	cmd := exec.Command(
		mainfestToolPath(),
		args...,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func mainfestToolPath() string {
	if runtime.GOOS == "windows" {
		return "C:/bin/manifest-tool.exe"
	}

	return "/bin/manifest-tool"
}
