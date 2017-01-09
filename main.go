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
	    git-release [--major|--minor|--patch] [--stable|--beta|--rc] [--force]
	    git-release --stable|--beta|--rc [--silent]
	    git-release <version> [--force] [--silent]
	    git-release --help
	    git-release --version

	  Options:
	    -h --help     Show this screen.
	    --version     Show version.
	    --dry-run
	    --force
	`)

	arguments, _ := docopt.Parse(usage, nil, true, "Git Release 0.1.0", false)

	silent := arguments["--silent"] == true
	force := arguments["--force"] == true

	if git.IsDirty() && !force {
		errors.Exit(errors.DirtyWorkspace, silent)
	}

	if git.CurrentTag() != "" && !force {
		errors.Exit(errors.AlreadyTagged, silent)
	}

	version, _ := arguments["<version>"].(string)

	if version == "" {
		lastTag := git.LastTag()
		if lastTag != "" {
			version = utils.NextVersion(git.LastTag(), arguments)
		} else {
			errors.Exit(errors.NoCurrentTagFound, silent)
		}
	}

	if version == "" {
		errors.Exit(errors.NotTagFound, silent)
	}

	if git.TagExists(version) {
		errors.Exit(errors.TagExists, silent)
	}

	versionfiles.UpdateAll(version)

	if git.IsDirty() {
		err := git.CommitAll("Release " + version)
		if err != nil {
			errors.Exit(errors.CommitFailed, silent)
		}

		fmt.Println("Commit: \"Release " + version + "\"")
	}

	git.Tag(version)
	fmt.Println("Tag:    \"" + version + "\"")
}
