package versionfiles

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

// UpdateAll :
func UpdateAll(version string) error {

	updateJSON("package.json", version)
	updateJSON("bower.json", version)
	updateJSON("composer.json", version)

	for _, name := range find("*.podspec.json") {
		updateJSON(name, version)
	}

	for _, name := range find("*.podspec") {
		updatePodspec(name, version)
	}

	if _, err := os.Stat("./.git-release/set-version.sh"); err == nil {
		cmd := exec.Command("./.git-release/set-version.sh", version)
		if err := cmd.Run(); err != nil {
			return err
		}

	}
	return nil
}

func updateJSON(name string, version string) {

	fileStat, err := os.Stat(name)
	if os.IsNotExist(err) {
		return
	}

	bytes, _ := ioutil.ReadFile(name)
	content := string(bytes)

	r, _ := regexp.Compile("\"version\"\\s*:\\s*\"[^\"]+\"")

	content = r.ReplaceAllString(content, "\"version\": \""+version+"\"")

	ioutil.WriteFile(name, []byte(content), fileStat.Mode())
}

func updatePodspec(name string, version string) {

	fileStat, err := os.Stat(name)
	if os.IsNotExist(err) {
		return
	}

	bytes, _ := ioutil.ReadFile(name)
	content := string(bytes)

	r, _ := regexp.Compile("(\\.version\\s*=\\s*)(['\"])[^'\"]+['\"]")

	content = r.ReplaceAllString(content, "${1}${2}"+version+"${2}")

	ioutil.WriteFile(name, []byte(content), fileStat.Mode())
}

func find(pattern string) []string {
	files, _ := filepath.Glob(pattern)
	return files
}
