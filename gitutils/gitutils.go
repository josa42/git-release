package gitutils

import (
	"os/exec"
	"strings"
)

var defaultFlags = []string{}

// var defaultFlags = []string{"--porcelain", "--no-color"}

// Exec :
func Exec(args ...string) (string, error) {

	cmd := gitCommand(args...)
	outputBytes, err := cmd.Output()

	// if err != nil {
	// 	fmt.Println(err)
	// }

	return strings.Trim(string(outputBytes), " \n"), err
}

// IsDirty :
func IsDirty() bool {

	out, _ := Exec("status", "--porcelain")

	return out != ""
}

// Tags :
func Tags() []string {
	out, _ := Exec("tag", "--list")

	if out != "" {
		// TODO sorting
		return strings.Split(out, "\n")
	}

	return []string{}
}

// AddAll :
func AddAll() {
	Exec("add", "--update")
}

// Commit :
func Commit(message string) error {
	_, error := Exec("commit", "--message", message)

	return error
}

// CommitAll :
func CommitAll(message string) error {
	AddAll()
	return Commit(message)
}

// Tag :
func Tag(version string) error {
	_, error := Exec("tag", version)

	return error
}

// LastTag :
func LastTag() string {
	out, _ := Exec("describe", "--tags", "--abbrev=0")
	return out
}

// CurrentTag :
func CurrentTag() string {
	out, _ := Exec("describe", "--tags", "--exact")
	return out
}

// TagExists :
func TagExists(tag string) bool {
	for _, existingTag := range Tags() {
		if existingTag == tag {
			return true
		}
	}

	return false
}

// // Tag :
// func TagExists(tag string) error {
//
// }

func gitCommand(args ...string) *exec.Cmd {

	// fmt.Println("git", args)

	cmd := exec.Command("git")

	for _, arg := range args {
		cmd.Args = append(cmd.Args, arg)
	}

	for _, globalArg := range defaultFlags {
		cmd.Args = append(cmd.Args, globalArg)
	}

	return cmd
}