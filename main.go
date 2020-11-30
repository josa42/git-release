package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/josa42/git-release/errors"
	"github.com/josa42/git-release/utils"
	"github.com/josa42/git-release/versionfiles"

	git "github.com/josa42/go-gitutils"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use: "git-release [<version>]",
	Example: strings.Join([]string{
		"  git-release v1.0.0",
		"  git-release --major",
	}, "\n"),

	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return fmt.Errorf("Too many arguments")
		}

		versionFlags := []string{"major", "minor", "patch"}
		modFlags := []string{"stable", "rc", "beta"}
		if len(args) == 1 {
			for _, f := range append(versionFlags, modFlags...) {
				if v, _ := cmd.Flags().GetBool(f); v {
					return fmt.Errorf("--%s cannot be used in combination with version", f)
				}
			}
		}

		if len(args) == 0 {

			found := false

			for _, f := range versionFlags {
				if v, _ := cmd.Flags().GetBool(f); v {
					found = true
				}
			}

			if !found {
				return fmt.Errorf("--%s or version needs to be provided", strings.Join(versionFlags, ", --"))
			}

		}

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {

		force, _ := cmd.Flags().GetBool("force")
		dirty, _ := cmd.Flags().GetBool("dirty")
		noPush, _ := cmd.Flags().GetBool("do-not-push")
		noEmptyCommit, _ := cmd.Flags().GetBool("no-empty-commit")
		messageTpl, _ := cmd.Flags().GetString("message")

		major, _ := cmd.Flags().GetBool("major")
		minor, _ := cmd.Flags().GetBool("minor")
		patch, _ := cmd.Flags().GetBool("patch")

		stable, _ := cmd.Flags().GetBool("stable")
		beta, _ := cmd.Flags().GetBool("beta")
		rc, _ := cmd.Flags().GetBool("rc")

		if !cmd.Flags().Lookup("message").Changed {
			out, err := exec.Command("git", "config", "git-release.message").Output()
			if err == nil {
				messageTpl = strings.TrimSpace(string(out))
			}
		}

		if !strings.Contains(messageTpl, "{version}") {
			errors.Exit(errors.InvalidMassage)
		}

		fmt.Printf("force         = %v\n", force)
		fmt.Printf("dirty         = %v\n", dirty)
		fmt.Printf("noPush        = %v\n", noPush)
		fmt.Printf("noEmptyCommit = %v\n", noEmptyCommit)
		fmt.Printf("messageTpl    = %v\n", messageTpl)

		if git.IsDirty() && !dirty {
			errors.Exit(errors.DirtyWorkspace)
		}

		if git.CurrentTag() != "" && noEmptyCommit && !force {
			errors.Exit(errors.AlreadyTagged)
		}

		vPrefix := false // TODO add option

		version := ""
		if len(args) == 1 {
			version = args[0]
		}

		if strings.HasPrefix(version, "v") {
			vPrefix = true
		}

		if version == "" {
			lastTag := git.LastTag()
			if lastTag != "" {
				// TODO still respect config
				if strings.HasPrefix(lastTag, "v") {
					vPrefix = true
				}
				version = utils.NextVersion(git.LastTag(), utils.VersionOptions{
					Major:  major,
					Minor:  minor,
					Patch:  patch,
					Stable: stable,
					Rc:     rc,
					Beta:   beta,
				})
			} else {
				errors.Exit(errors.NoCurrentTagFound)
			}
		}

		if version == "" {
			errors.Exit(errors.NotTagFound)
		}

		if vPrefix && !strings.HasPrefix(version, "v") {
			version = fmt.Sprintf("v%s", version)
		}

		fmt.Printf("version = %v\n", version)

		if git.TagExists(version) {
			errors.Exit(errors.TagExists)
		}

		if err := versionfiles.UpdateAll(version); err != nil {
			errors.Exit(errors.UpdateVersionFailed)
		}

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

		if !noPush {
			if err := git.Push(); err != nil {
				errors.Exit(errors.PushFailed)
			}
		}
	},
}

func init() {
	cmd.Flags().BoolP("major", "", false, "")
	cmd.Flags().BoolP("minor", "", false, "")
	cmd.Flags().BoolP("patch", "", false, "")

	cmd.Flags().BoolP("stable", "", false, "")
	cmd.Flags().BoolP("rc", "", false, "")
	cmd.Flags().BoolP("beta", "", false, "")

	cmd.Flags().StringP("message", "", "Release {version}", "Commit message")
	cmd.Flags().BoolP("version", "", false, "Show version")
	cmd.Flags().BoolP("dirty", "", false, "Include changed files in release commit")
	cmd.Flags().BoolP("force", "", false, "Force new commit even thought the latest commit is already tagged")
	cmd.Flags().BoolP("do-not-push", "", false, "Do not push commit and tags")
	cmd.Flags().BoolP("no-empty-commit", "", false, "Do not commit if nothing changed")
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
