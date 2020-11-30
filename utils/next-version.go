package utils

import (
	"strings"

	"github.com/blang/semver"
)

type VersionOptions struct {
	Major  bool
	Minor  bool
	Patch  bool
	Stable bool
	Rc     bool
	Beta   bool
}

// NextVersion :
func NextVersion(version string, options VersionOptions) string {

	version = strings.TrimPrefix(version, "v")

	next, err := semver.Parse(version)
	if err != nil {
		panic(err)
	}

	current, _ := semver.Parse(version)

	if options.Major {
		next.Major++
		next.Minor = 0
		next.Patch = 0
		resetPreVersion(&next)

	} else if options.Minor {
		next.Minor++
		next.Patch = 0
		resetPreVersion(&next)

	} else if options.Patch {
		next.Patch++
		resetPreVersion(&next)
	}

	if options.Rc {
		increasePreVersion("rc", &next)

	} else if options.Beta {
		increasePreVersion("beta", &next)

	} else if options.Stable {
		next.Pre = []semver.PRVersion{}
	}

	if next.GT(current) {
		return next.String()
	}

	panic("Invalid version step")
}

func resetPreVersion(version *semver.Version) {
	version.Pre = []semver.PRVersion{}
}

func increasePreVersion(preType string, version *semver.Version) {

	preLen := len(version.Pre)

	preVer, _ := semver.NewPRVersion(preType)

	if preLen > 0 {
		currentPreVer := &version.Pre[0]

		if currentPreVer.Compare(preVer) == 0 {
			preNum, _ := semver.NewPRVersion("2")

			if preLen == 2 {
				currentPreNum := &version.Pre[1]

				if currentPreNum.IsNumeric() {
					preNum.VersionNum = currentPreNum.VersionNum + 1
				}
			}

			version.Pre = []semver.PRVersion{preVer, preNum}
		} else {
			version.Pre = []semver.PRVersion{preVer}
		}
	} else {
		version.Pre = []semver.PRVersion{preVer}
	}
}
