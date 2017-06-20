package main

import (
	"fmt"
	"strings"

	docopt "github.com/docopt/docopt-go"
	"github.com/josa42/git-release/errors"
	"github.com/josa42/git-release/utils"
	"github.com/josa42/git-release/versionfiles"
	git "github.com/josa42/go-gitutils"
	stringutils "github.com/josa42/go-stringutils"
)

func main() {
	usage := stringutils.TrimLeadingTabs(`
		Usage:
		  git-release [--major|--minor|--patch] [--stable|--beta|--rc] [--dirty] [--force] [--do-not-push] [--no-empty-commit] [--message=<message>]
		  git-release --stable|--beta|--rc                             [--dirty] [--force] [--do-not-push] [--no-empty-commit] [--message=<message>]
		  git-release <version>                                        [--dirty] [--force] [--do-not-push] [--no-empty-commit] [--message=<message>]
		  git-release --help
		  git-release --version

		Options:
		  -h --help          Show this screen.
		  --version          Show version.
		  --dirty            Include changed files in release commit.
		  --force            Force new commit even thought the latest commit is already tagged.
		  --do-not-push      Do not push commit and tags
		  --no-empty-commit  Do not commit if nothing changed
  `)

	arguments, _ := docopt.Parse(usage, nil, true, "Git Release 0.1.0", false)

	force := arguments["--force"] == true
	dirty := arguments["--dirty"] == true
	push := arguments["--no-empty-commit"] != true
	noEmptyCommit := arguments["--no-empty-commit"] == true
	messageTpl := "🎉  Release {version}"
	if arguments["--message"] != nil {
		messageTpl = arguments["--message"].(string)
	}

	if !strings.Contains(messageTpl, "{version}") {
		errors.Exit(errors.InvalidMassage)
	}

	if git.IsDirty() && !dirty {
		errors.Exit(errors.DirtyWorkspace)
	}

	if git.CurrentTag() != "" && noEmptyCommit && !force {
		errors.Exit(errors.AlreadyTagged)
	}

	version, _ := arguments["<version>"].(string)

	if version == "" {
		lastTag := git.LastTag()
		if lastTag != "" {
			version = utils.NextVersion(git.LastTag(), arguments)
		} else {
			errors.Exit(errors.NoCurrentTagFound)
		}
	}

	if version == "" {
		errors.Exit(errors.NotTagFound)
	}

	if git.TagExists(version) {
		errors.Exit(errors.TagExists)
	}

	versionfiles.UpdateAll(version)

	message := strings.Replace(messageTpl, "{version}", version, 1)

	if git.IsDirty() {
		err := git.CommitAll(message)
		if err != nil {
			errors.Exit(errors.CommitFailed)
		}

		fmt.Println("Commit: \"" + message + "\"")
	} else if !noEmptyCommit {
		err := git.CommitEmpty(message)
		if err != nil {
			errors.Exit(errors.CommitFailed)
		}

		fmt.Println("Commit: \"" + message + "\"")
	}

	git.Tag(version)
	fmt.Println("Tag:    \"" + version + "\"")

	if push {
		git.Push()
	}
}
