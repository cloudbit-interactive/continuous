package controller

import (
	"io/ioutil"
	"log"
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

func CreateFile(file string, content string) string {
	dirArray := strings.Split(file, "/")
	dirArray = dirArray[0 : len(dirArray)-1]
	dirString := strings.Join(dirArray, "/")
	CreateDir(dirString)
	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}
	Log("-- File: " + file)
	if content == "nil" {
		return file
	}
	content = strings.TrimSpace(ReplaceString(content))
	_, err2 := f.WriteString(content)
	if err2 != nil {
		log.Fatal(err2)
	}
	return file
}

func CreateDir(path string) string {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
	}
	return path
}

func DeletePath(path string) string {
	err := os.RemoveAll(path)
	if err != nil {
	}
	return path
}

func ExistPath(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func MovePath(from string, to string) string {
	fromFile, err := os.Open(from)
	if err != nil {

	}
	defer fromFile.Close()
	fileInfo, err := fromFile.Stat()
	if err != nil {

	}
	if fileInfo.IsDir() {
		CreateDir(to)
	} else {
		CreateFile(to, "")
		err := os.Rename(from, to)
		if err != nil {
			cuppago.Log(err)
		}
	}
	DeletePath(from)
	return to
}
