// Copyright (c) 2023, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

// Pipeline provides Pipeline metadata from the environment.
type Pipeline struct {
	// Build provides build metadata.
	Build struct {
		Branch   string `envconfig:"DRONE_BUILD_BRANCH"`
		Number   int    `envconfig:"DRONE_BUILD_NUMBER"`
		Parent   int    `envconfig:"DRONE_BUILD_PARENT"`
		Event    string `envconfig:"DRONE_BUILD_EVENT"`
		Action   string `envconfig:"DRONE_BUILD_ACTION"`
		Status   string `envconfig:"DRONE_BUILD_STATUS"`
		Created  int64  `envconfig:"DRONE_BUILD_CREATED"`
		Started  int64  `envconfig:"DRONE_BUILD_STARTED"`
		Finished int64  `envconfig:"DRONE_BUILD_FINISHED"`
		Link     string `envconfig:"DRONE_BUILD_LINK"`
	}

	// Calver provides the calver details parsed from the
	// git tag. If the git tag is empty or is not a valid
	// calver, the values will be empty.
	Calver struct {
		Version    string `envconfig:"DRONE_CALVER"`
		Short      string `envconfig:"DRONE_CALVER_SHORT"`
		MajorMinor string `envconfig:"DRONE_CALVER_MAJOR_MINOR"`
		Major      string `envconfig:"DRONE_CALVER_MAJOR"`
		Minor      string `envconfig:"DRONE_CALVER_MINOR"`
		Micro      string `envconfig:"DRONE_CALVER_MICRO"`
		Modifier   string `envconfig:"DRONE_CALVER_MODIFIER"`
	}

	// Card provides adaptive card configuration options.
	Card struct {
		Path string `envconfig:"DRONE_CARD_PATH"`
	}

	// Commit provides the commit metadata.
	Commit struct {
		Rev     string `envconfig:"DRONE_COMMIT_SHA"`
		Before  string `envconfig:"DRONE_COMMIT_BEFORE"`
		After   string `envconfig:"DRONE_COMMIT_AFTER"`
		Ref     string `envconfig:"DRONE_COMMIT_REF"`
		Branch  string `envconfig:"DRONE_COMMIT_BRANCH"`
		Source  string `envconfig:"DRONE_COMMIT_SOURCE"`
		Target  string `envconfig:"DRONE_COMMIT_TARGET"`
		Link    string `envconfig:"DRONE_COMMIT_LINK"`
		Message string `envconfig:"DRONE_COMMIT_MESSAGE"`

		Author struct {
			Username string `envconfig:"DRONE_COMMIT_AUTHOR"`
			Name     string `envconfig:"DRONE_COMMIT_AUTHOR_NAME"`
			Email    string `envconfig:"DRONE_COMMIT_AUTHOR_EMAIL"`
			Avatar   string `envconfig:"DRONE_COMMIT_AUTHOR_AVATAR"`
		}
	}

	// Deploy provides the deployment metadata.
	Deploy struct {
		ID     string `envconfig:"DRONE_DEPLOY_TO"`
		Target string `envconfig:"DRONE_DEPLOY_ID"`
	}

	// Failed provides a list of failed steps and failed stages
	// for the current pipeline.
	Failed struct {
		Steps  []string `envconfig:"DRONE_FAILED_STEPS"`
		Stages []string `envconfig:"DRONE_FAILED_STAGES"`
	}

	// Git provides the git repository metadata.
	Git struct {
		HTTPURL string `envconfig:"DRONE_GIT_HTTP_URL"`
		SSHURL  string `envconfig:"DRONE_GIT_SSH_URL"`
	}

	// PullRequest provides the pull request metadata.
	PullRequest struct {
		Number int `envconfig:"DRONE_PULL_REQUEST"`
	}

	// Repo provides the repository metadata.
	Repo struct {
		Branch     string `envconfig:"DRONE_REPO_BRANCH"`
		Link       string `envconfig:"DRONE_REPO_LINK"`
		Namespace  string `envconfig:"DRONE_REPO_NAMESPACE"`
		Name       string `envconfig:"DRONE_REPO_NAME"`
		Private    bool   `envconfig:"DRONE_REPO_PRIVATE"`
		Remote     string `envconfig:"DRONE_GIT_HTTP_URL"`
		SCM        string `envconfig:"DRONE_REPO_SCM"`
		Slug       string `envconfig:"DRONE_REPO"`
		Visibility string `envconfig:"DRONE_REPO_VISIBILITY"`
	}

	// Stage provides the stage metadata.
	Stage struct {
		Kind      string   `envconfig:"DRONE_STAGE_KIND"`
		Type      string   `envconfig:"DRONE_STAGE_TYPE"`
		Name      string   `envconfig:"DRONE_STAGE_NAME"`
		Number    int      `envconfig:"DRONE_STAGE_NUMBER"`
		Machine   string   `envconfig:"DRONE_STAGE_MACHINE"`
		OS        string   `envconfig:"DRONE_STAGE_OS"`
		Arch      string   `envconfig:"DRONE_STAGE_ARCH"`
		Variant   string   `envconfig:"DRONE_STAGE_VARIANT"`
		Status    string   `envconfig:"DRONE_STAGE_STATUS"`
		Started   int64    `envconfig:"DRONE_STAGE_STARTED"`
		Finished  int64    `envconfig:"DRONE_STAGE_FINISHED"`
		DependsOn []string `envconfig:"DRONE_STAGE_DEPENDS_ON"`
	}

	// Step provides the step metadata.
	Step struct {
		Number int    `envconfig:"DRONE_STEP_NUMBER"`
		Name   string `envconfig:"DRONE_STEP_NAME"`
	}

	// Semver provides the semver details parsed from the
	// git tag. If the git tag is empty or is not a valid
	// semver, the values will be empty and the error field
	// will be populated with the parsing error.
	Semver struct {
		Version    string `envconfig:"DRONE_SEMVER"`
		Short      string `envconfig:"DRONE_SEMVER_SHORT"`
		Major      string `envconfig:"DRONE_SEMVER_MAJOR"`
		Minor      string `envconfig:"DRONE_SEMVER_MINOR"`
		Patch      string `envconfig:"DRONE_SEMVER_PATCH"`
		Build      string `envconfig:"DRONE_SEMVER_BUILD"`
		PreRelease string `envconfig:"DRONE_SEMVER_PRERELEASE"`
		Error      string `envconfig:"DRONE_SEMVER_ERROR"`
	}

	// System provides the Drone system metadata, including
	// the system version of details required to create the
	// drone website address.
	System struct {
		Proto    string `envconfig:"DRONE_SYSTEM_PROTO"`
		Host     string `envconfig:"DRONE_SYSTEM_HOST"`
		Hostname string `envconfig:"DRONE_SYSTEM_HOSTNAME"`
		Version  string `envconfig:"DRONE_SYSTEM_VERSION"`
	}

	// Tag provides the git tag details.
	Tag struct {
		Name string `envconfig:"DRONE_TAG"`
	}
}
