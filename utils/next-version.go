package utils

import "github.com/blang/semver"

// NextVersion :
func NextVersion(version string, options map[string]interface{}) string {

	next, err := semver.Parse(version)
	if err != nil {
		panic(err)
	}

	current, _ := semver.Parse(version)

	if options["--major"] == true {
		next.Major++
		next.Minor = 0
		next.Patch = 0
		resetPreVersion(&next)

	} else if options["--minor"] == true {
		next.Minor++
		next.Patch = 0
		resetPreVersion(&next)

	} else if options["--patch"] == true {
		next.Patch++
		resetPreVersion(&next)
	}

	if options["--rc"] == true {
		increasePreVersion("rc", &next)

	} else if options["--beta"] == true {
		increasePreVersion("beta", &next)

	} else if options["--stable"] == true {
		next.Pre = []semver.PRVersion{}
	}

	if next.GT(current) {
		return next.String()
	}

	panic("Invalid version step")

	return ""
	// panic("Invalid version step")
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
