package controller

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/cloudbit-interactive/cuppago"
)

func AutoGetFile() string {
	files, err := ioutil.ReadDir(cuppago.GetRootPath())
	if err != nil {
		cuppago.Error(err)
	}
	//_ := ""
	for _, file := range files {
		name := file.Name()
		if strings.Contains(name, ".yaml") == true {
			return name
		}
	}
	return ""
}

func GetFileContent(file string) string {
	text, err := os.ReadFile(file)
	if err != nil {
		cuppago.Error(err)
	}
	return string(text)
}
