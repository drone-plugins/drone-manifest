// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/LemontechSA/drone-manifest-ecr/tagging"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/drone/drone-template-lib/template"
	"github.com/urfave/cli/v2"
)

// Settings for the plugin.
type Settings struct {
	AccessKey     string
	SecretKey     string
	Region        string
	AssumeRole    string
	ExternalID    string
	Insecure      bool
	Platforms     cli.StringSlice
	Target        string
	Template      string
	Spec          string
	IgnoreMissing bool
	Tags          cli.StringSlice
	AutoTag       bool
}

const defaultRegion = "us-east-1"

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	if p.settings.AccessKey == "" && p.settings.AccessKey != "" {
		return errors.New("you must provide a username")
	}

	if p.settings.SecretKey == "" && p.settings.SecretKey != "" {
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
			t.Build.Tags = tagging.DefaultTags(p.pipeline.Commit.Ref, p.pipeline.Commit.SHA)
		} else {
			log.Printf("skipping automated tags for %s", p.pipeline.Commit.Ref)
			return nil
		}
	} else {
		t.Build.Tags = p.settings.Tags.Value()
	}

	region := p.settings.Region
	assumeRole := p.settings.AssumeRole
	externalID := p.settings.ExternalID

	if region == "" {
		region = defaultRegion
	}

	key := p.settings.AccessKey
	secret := p.settings.SecretKey

	if key != "" && secret != "" {
		os.Setenv("AWS_ACCESS_KEY_ID", key)
		os.Setenv("AWS_SECRET_ACCESS_KEY", secret)
	}

	sess, err := session.NewSession(&aws.Config{Region: &region})
	if err != nil {
		log.Fatalf(fmt.Sprintf("error creating aws session: %v", err))
	}

	svc := getECRClient(sess, assumeRole, externalID)
	username, password, _, err := getAuthInfo(svc)

	if err != nil {
		log.Fatalf(fmt.Sprintf("error getting ECR auth: %v", err))
	}

	args := []string{
		fmt.Sprintf("--username=%s", username),
		fmt.Sprintf("--password=%s", password),
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
			fmt.Sprintf("--template=%s", p.settings.Template),
		)

		if !p.settings.AutoTag {
			args = append(args, fmt.Sprintf("--target=%s", p.settings.Target))
		}
	}

	if p.settings.IgnoreMissing {
		args = append(args, "--ignore-missing")
	}

	for _, tag := range t.Build.Tags {
		args = append(args, fmt.Sprintf("--target=%s", p.settings.Target+":"+tag))

		cmd := exec.Command(
			mainfestToolPath(),
			args...,
		)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()

		if err != nil {
			return err
		}
	}

	cmd := exec.Command(
		mainfestToolPath(),
		args...,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func getECRClient(sess *session.Session, role string, externalID string) *ecr.ECR {
	if role == "" {
		return ecr.New(sess)
	}
	if externalID != "" {
		return ecr.New(sess, &aws.Config{
			Credentials: stscreds.NewCredentials(sess, role, func(p *stscreds.AssumeRoleProvider) {
				p.ExternalID = &externalID
			}),
		})
	}

	return ecr.New(sess, &aws.Config{
		Credentials: stscreds.NewCredentials(sess, role),
	})
}

func getAuthInfo(svc *ecr.ECR) (username, password, registry string, err error) {
	var result *ecr.GetAuthorizationTokenOutput
	var decoded []byte

	result, err = svc.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return
	}

	auth := result.AuthorizationData[0]
	token := *auth.AuthorizationToken
	decoded, err = base64.StdEncoding.DecodeString(token)
	if err != nil {
		return
	}

	registry = strings.TrimPrefix(*auth.ProxyEndpoint, "https://")
	creds := strings.Split(string(decoded), ":")
	username = creds[0]
	password = creds[1]
	return
}

func mainfestToolPath() string {
	if runtime.GOOS == "windows" {
		return "C:/bin/manifest-tool.exe"
	}

	return "/bin/manifest-tool"
}
