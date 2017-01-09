package main

import (
	"fmt"

	docopt "github.com/docopt/docopt-go"
	"github.com/josa42/git-release/errors"
	git "github.com/josa42/git-release/gitutils"
	"github.com/josa42/git-release/stringutils"
	"github.com/josa42/git-release/utils"
	"github.com/josa42/git-release/versionfiles"
)

func main() {
	usage := stringutils.StripPrefix(`
	  Usage:
	    git-release [--major|--minor|--patch] [--stable|--beta|--rc] [--dirty] [--force] [--no-empty-commit]
	    git-release --stable|--beta|--rc [--dirty] [--force] [--no-empty-commit]
	    git-release <version> [--dirty] [--force] [--no-empty-commit]
	    git-release --help
	    git-release --version

	  Options:
	    -h --help          Show this screen.
	    --version          Show version.
	    --dirty            Include changed files in release commit.
	    --force            Force new commit even thought the latest commit is already tagged.
      --no-empty-commit  Do not commit if nothing changed
	`)

	arguments, _ := docopt.Parse(usage, nil, true, "Git Release 0.1.0", false)

	force := arguments["--force"] == true
	dirty := arguments["--dirty"] == true
	push := arguments["--push"] == true
	noEmptyCommit := arguments["--no-empty-commit"] == true

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

	if git.IsDirty() {
		err := git.CommitAll("Release " + version)
		if err != nil {
			errors.Exit(errors.CommitFailed)
		}

		fmt.Println("Commit: \"Release " + version + "\"")
	} else if !noEmptyCommit {
		err := git.CommitEmpty("Release " + version)
		if err != nil {
			errors.Exit(errors.CommitFailed)
		}

		fmt.Println("Commit: \"Release " + version + "\"")
	}

	git.Tag(version)
	fmt.Println("Tag:    \"" + version + "\"")

	if push {
		git.Push()
	}
}
