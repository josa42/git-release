package errors

import (
	"fmt"
	"os"
)

// NotImplemented :
const NotImplemented = 2

// DirtyWorkspace :
const DirtyWorkspace = 3

// TagExists :
const TagExists = 4

// NotTagFound :
const NotTagFound = 5

// NoCurrentTagFound :
const NoCurrentTagFound = 6

// CommitFailed :
const CommitFailed = 7

// AlreadyTagged :
const AlreadyTagged = 8

var messages = map[int]string{
	TagExists:         "Tag exists.",
	NotTagFound:       "Could not determ tag.",
	NoCurrentTagFound: "Could not find current tag.",
	NotImplemented:    "Not implemented yet.",
	DirtyWorkspace:    "Dirty workspace. Use --dirty to include changed files in the release commit.",
	CommitFailed:      "Commit failed.",
	AlreadyTagged:     "Commit is already tagged. Use --force to allow multiple tags.",
}

// Exit :
func Exit(code int) {

	if messages[code] != "" {
		fmt.Println(messages[code])
	}

	os.Exit(code)
}
