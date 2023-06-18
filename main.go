// Copyright (c) 2023, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package main

import (
	"bytes"
	"context"
	"strings"

	"github.com/drone-plugins/drone-manifest/plugin"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(formatter))

	var args plugin.Args
	if err := envconfig.Process("", &args); err != nil {
		logrus.Fatalln(err)
	}

	switch args.Level {
	case "debug":
		logrus.SetFormatter(textFormatter)
		logrus.SetLevel(logrus.DebugLevel)
	case "trace":
		logrus.SetFormatter(textFormatter)
		logrus.SetLevel(logrus.TraceLevel)
	}

	if err := plugin.Exec(context.Background(), &args); err != nil {
		logrus.Fatalln(err)
	}
}

// default formatter that writes logs without including timestamp or level information.
type formatter struct{}

func (*formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	b.WriteString(entry.Message)
	if !strings.HasSuffix(entry.Message, "\n") {
		b.WriteByte('\n')
	}

	return b.Bytes(), nil
}

// text formatter that writes logs with level information.
var textFormatter = &logrus.TextFormatter{
	DisableTimestamp: true,
}
