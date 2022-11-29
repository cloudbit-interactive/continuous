package controller

import (
	"bytes"
	"fmt"
	"github.com/cloudbit-interactive/cuppago"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func ReplaceVariables(values map[string]interface{}) map[string]interface{} {
	for key := range values {
		values[key] = ReplaceToSystemValue(cuppago.String(values[key]))
	}
	return values
}

func ReplaceString(string string) string {
	for key := range YamlVars {
		string = cuppago.ReplaceNotCase(string, "\\${"+key+"}", cuppago.String(YamlVars[key]))
	}
	string = ReplaceToSystemValue(string)
	return string
}

func ReplaceToSystemValue(string string) string {
	os := runtime.GOOS
	if os == "darwin" {
		os = "mac"
	}
	arch := runtime.GOARCH
	date := time.Now().String()
	string = cuppago.ReplaceNotCase(string, "\\${DATE}", date[0:10])
	string = cuppago.ReplaceNotCase(string, "\\${DATETIME}", date[0:19])
	string = cuppago.ReplaceNotCase(string, "\\${OS}", os)
	string = cuppago.ReplaceNotCase(string, "\\${ARCH}", arch)
	string = cuppago.ReplaceNotCase(string, "\\${ARCH}", arch)
	if len(YamlOutput) != 0 {
		output := fmt.Sprint(YamlOutput[len(YamlOutput)-1])
		string = cuppago.ReplaceNotCase(string, "\\${OUTPUT}", output)
	}
	return string
}

func BashCommand(command string) string {
	outputString := ""
	command = ReplaceString(command)
	Log("-- CMD: " + command)

	if runtime.GOOS == "windows" {
		output, err := exec.Command("powershell", command).Output()
		if err != nil {
			outputString = cuppago.String(err)
		} else {
			outputString = string(output)
		}
	} else {
		output, err := exec.Command("bash", "-c", command).Output()
		if err != nil {
			outputString = cuppago.String(err)
		} else {
			outputString = string(output)
		}
	}
	outputString = strings.TrimSpace(outputString)
	outputString = strings.Trim(outputString, "\n")
	Log("---- output: " + outputString)
	return outputString
}

func Command(app string, args []string, workingDirectory string, backgroundName string) string {
	for i := 0; i < len(args); i++ {
		args[i] = strings.TrimSpace(ReplaceString(args[i]))
	}
	workingDirectory = strings.TrimSpace(ReplaceString(workingDirectory))
	Log("-- CMD: "+app, "-- args: ", args, "-- workingDirectory: "+workingDirectory)
	var output bytes.Buffer
	cmd := exec.Command(app, args...)
	if backgroundName != "" {
		YamlBackgrounds[backgroundName] = cmd
	}
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
	Log("---- output: " + outputString)
	return outputString
}

func Log(values ...interface{}) {
	if YamlData["log"] != true {
		return
	}
	cuppago.LogFile(values...)
}

func KillPorts(port string) string {
	ports := strings.Split(port, ",")
	output := ""
	for _, value := range ports {
		output = KillPort(value)
	}
	return output
}

func KillPort(port string) string {
	port = strings.TrimSpace(port)
	output := ""
	if runtime.GOOS == "windows" {
		args := []string{"Stop-Process", "-Id", fmt.Sprintf("(Get-NetTCPConnection -LocalPort %s).OwningProcess", port)}
		output = Command("powershell", args, "", "")
	} else {
		command := fmt.Sprintf("lsof -i tcp:%s | grep LISTEN | awk '{print $2}' | xargs kill -9", port)
		output = BashCommand(command)
	}
	return output
}
