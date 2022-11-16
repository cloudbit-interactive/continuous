package controller

import (
	"bytes"
	"github.com/cloudbit-interactive/cuppago"
	"os/exec"
	"strings"
)

func ReplaceString(string string) string {
	for key := range YamlVars {
		string = cuppago.ReplaceNotCase(string, "\\${"+key+"}", cuppago.String(YamlVars[key]))
	}
	return string
}

func BashCommand(command string) string {
	outputString := ""
	command = ReplaceString(command)
	output, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		outputString = cuppago.String(err)
	} else {
		outputString = string(output)
	}
	outputString = strings.TrimSpace(outputString)
	outputString = strings.Trim(outputString, "\n")
	cuppago.LogFile("CMD", command, outputString)
	return outputString
}

func Command(app string, args []string, workingDirectory string) string {
	for i := 0; i < len(args); i++ {
		args[i] = strings.TrimSpace(ReplaceString(args[i]))
	}
	workingDirectory = strings.TrimSpace(ReplaceString(workingDirectory))
	var output bytes.Buffer
	cmd := exec.Command(app, args...)
	cmd.Dir = workingDirectory
	cmd.Stdout = &output
	err := cmd.Run()
	outputString := ""
	if err != nil {
		outputString = cuppago.String(err)
	} else {
		outputString = output.String()
	}
	outputString = strings.TrimSpace(outputString)
	outputString = strings.Trim(outputString, "\n")
	cuppago.LogFile("CMD:", app, args, workingDirectory, outputString)
	return outputString
}
