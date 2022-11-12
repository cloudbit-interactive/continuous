package controller

import (
	"os"

	"github.com/cloudbit-interactive/cuppago"
)

func GetFileContent(file string) string {
	text, err := os.ReadFile(file)
	if err != nil {
		cuppago.Error(err)
	}
	return string(text)
}
