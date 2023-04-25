// Copyright (c) 2023, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"os"
	"strconv"
)

type (
	legacyRepo struct {
		Owner  string
		Name   string
		Branch string
	}

	legacyBuild struct {
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

	legacyJob struct {
		Started int64
	}

	legacyPlugin struct {
		Repo  legacyRepo
		Build legacyBuild
		Job   legacyJob
	}
)

func toLegacyPlugin(args *Args) legacyPlugin {
	return legacyPlugin{
		Repo: legacyRepo{
			Owner:  args.Repo.Namespace,
			Name:   args.Repo.Name,
			Branch: args.Repo.Branch,
		},
		Build: legacyBuild{
			Path:    getEnv("DRONE_WORKSPACE", ""),
			Tag:     args.Tag.Name,
			Number:  args.Build.Number,
			Event:   args.Build.Event,
			Status:  args.Build.Status,
			Commit:  args.Commit.Rev,
			Ref:     args.Commit.Ref,
			Branch:  args.Commit.Branch,
			Pull:    strconv.FormatInt(int64(args.PullRequest.Number), 10), // c.String("commit.pull"),
			Started: args.Build.Started,
			Created: args.Build.Created,
			Tags:    args.Tags,
		},
		Job: legacyJob{
			Started: 0,
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
