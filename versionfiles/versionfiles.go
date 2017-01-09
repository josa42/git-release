package versionfiles

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

// UpdateAll :
func UpdateAll(version string) {

	updateJSON("package.json", version)
	updateJSON("bower.json", version)
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

	fmt.Println(name)
	fmt.Print(content)
}
