package command

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

type Command struct {
	username      string
	password      string
	spec          string
	platforms     []string
	target        string
	template      string
	path          string
	ignoreMissing bool
}

func New(opts ...Option) *Command {
	c := &Command{}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Command) Exec() error {
	args := []string{}

	if c.username != "" {
		args = append(args, fmt.Sprintf("--username=%s", c.username))
	}

	if c.password != "" {
		args = append(args, fmt.Sprintf("--password=%s", c.password))
	}

	args = append(args, "push")

	if c.spec != "" {
		tmpfile, err := ioutil.TempFile(c.path, "manifest-")

		if err != nil {
			return errors.Wrap(err, "failed to create tempfile")
		}

		defer os.Remove(tmpfile.Name())

		if _, err := tmpfile.Write([]byte(c.spec)); err != nil {
			return errors.Wrap(err, "failed to write tempfile")
		}

		if err := tmpfile.Close(); err != nil {
			return errors.Wrap(err, "failed to close temp file")
		}

		args = append(args, "from-spec")
		args = append(args, tmpfile.Name())
	} else {
		args = append(args, "from-args")

		if len(c.platforms) != 0 {
			args = append(args, fmt.Sprintf("--platforms=%s", strings.Join(c.platforms, ",")))
		}

		if c.target != "" {
			args = append(args, fmt.Sprintf("--target=%s", c.target))
		}

		if c.template != "" {
			args = append(args, fmt.Sprintf("--template=%s", c.template))
		}
	}

	if c.ignoreMissing {
		args = append(args, "--ignore-missing")
	}

	cmd := exec.Command(
		"manifest-tool",
		args...,
	)

	buf := bytes.NewBufferString("")
	cmd.Stdout = io.MultiWriter(os.Stdout, buf)
	cmd.Stderr = io.MultiWriter(os.Stderr, buf)

	if c.path != "" {
		cmd.Dir = c.path
	}

	return cmd.Run()
}
