// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package main

import (
	"github.com/drone-plugins/drone-manifest/plugin"
	"github.com/urfave/cli/v2"
)

// settingsFlags has the cli.Flags for the plugin.Settings.
func settingsFlags(settings *plugin.Settings) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "username",
			Usage:       "username for registry",
			EnvVars:     []string{"PLUGIN_USERNAME", "MANIFEST_USERNAME", "DOCKER_USERNAME"},
			Destination: &settings.Username,
		},
		&cli.StringFlag{
			Name:        "password",
			Usage:       "password for registry",
			EnvVars:     []string{"PLUGIN_PASSWORD", "MANIFEST_PASSWORD", "DOCKER_PASSWORD"},
			Destination: &settings.Password,
		},
		&cli.BoolFlag{
			Name:        "insecure",
			Usage:       "enable allow insecure registry",
			EnvVars:     []string{"PLUGIN_INSECURE"},
			Destination: &settings.Insecure,
		},
		&cli.StringSliceFlag{
			Name:        "platforms",
			Usage:       "platforms for manifests",
			EnvVars:     []string{"PLUGIN_PLATFORMS"},
			Destination: &settings.Platforms,
		},
		&cli.StringFlag{
			Name:        "target",
			Usage:       "target for manifests",
			EnvVars:     []string{"PLUGIN_TARGET"},
			Destination: &settings.Target,
		},
		&cli.StringFlag{
			Name:        "template",
			Usage:       "template for manifests",
			EnvVars:     []string{"PLUGIN_TEMPLATE"},
			Destination: &settings.Template,
		},
		&cli.StringFlag{
			Name:        "spec",
			Usage:       "path to manifest spec",
			EnvVars:     []string{"PLUGIN_SPEC"},
			Destination: &settings.Spec,
		},
		&cli.BoolFlag{
			Name:        "ignore-missing",
			Usage:       "ignore missing images",
			EnvVars:     []string{"PLUGIN_IGNORE_MISSING"},
			Destination: &settings.IgnoreMissing,
		},
		&cli.StringSliceFlag{
			Name:        "tags",
			Usage:       "list of additional tags",
			EnvVars:     []string{"PLUGIN_TAG", "PLUGIN_TAGS"},
			FilePath:    ".tags",
			Destination: &settings.Tags,
		},
		&cli.BoolFlag{
			Name:        "tags.auto",
			Usage:       "automatically build tags",
			EnvVars:     []string{"PLUGIN_DEFAULT_TAGS", "PLUGIN_AUTO_TAG"},
			Destination: &settings.AutoTag,
		},
	}
}
