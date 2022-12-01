package controller

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/cloudbit-interactive/cuppago"
	"gopkg.in/yaml.v3"
)

var YamlData map[string]interface{}
var YamlVars map[string]interface{}
var YamlBackgrounds map[string]interface{}
var YamlOutput []string

func ProcessYamlString(yamlString string, yamlVarsPath string) {
	YamlOutput = []string{}
	yamlData := make(map[string]interface{})
	YamlBackgrounds = make(map[string]interface{})
	err := yaml.Unmarshal([]byte(yamlString), yamlData)
	if err != nil {
		cuppago.Error(err)
	}
	YamlData = yamlData

	if yamlVarsPath != "" {
		YamlProcessVars(yamlVarsPath)
	} else {
		YamlProcessVars(yamlData["vars"])
	}

	jobs := YamlData["jobs"].([]interface{})
	if jobs == nil {
		Log("No jobs founds")
		return
	}
	if yamlData["port"] != nil {
		port := cuppago.String(yamlData["port"])
		cuppago.LogFile("Continuous running in http://localhost:" + port)
		YamlProcessJobs(jobs)
		http.Handle("/favicon.ico", http.NotFoundHandler())
		http.ListenAndServe(":"+port, nil)
	} else {
		cuppago.LogFile("Process running, press [Enter] to exit...")
		YamlProcessJobs(jobs)
		for {
			time.Sleep(time.Duration(1<<63 - 1))
		}
	}
}

func YamlProcessVars(vars interface{}) {
	if reflect.TypeOf(vars).String() == "string" {
		filePath := cuppago.GetRootPath() + "/" + vars.(string)
		text := GetFileContent(filePath)
		yamlData := make(map[string]interface{})
		err := yaml.Unmarshal([]byte(text), yamlData)
		if err != nil {
			cuppago.Error(err)
		}
		YamlVars = ReplaceVariables(yamlData)
	} else {
		YamlVars = ReplaceVariables(vars.(map[string]interface{}))
	}
}

func YamlProcessJobs(jobs []interface{}) {
	for i := 0; i < len(jobs); i++ {
		YamlJob(jobs[i].(map[string]interface{}))
	}
}

func YamlJob(job map[string]interface{}) {
	for key := range job {
		Log("JOB -----> " + key)
		if strings.Contains(key, "var-") == true {
			name := key[4:]
			YamlVars[name] = ReplaceString(job[key].(string))
		} else if key == CreateFileConst {
			output := YamlCreateFile(job[CreateFileConst].(map[string]interface{}))
			YamlOutput = append(YamlOutput, output)
		} else if key == CreateFolder {
			output := CreateDir(ReplaceString(job[CreateFolder].(string)))
			YamlOutput = append(YamlOutput, output)
		} else if key == Delete {
			output := DeletePath(ReplaceString(job[Delete].(string)))
			YamlOutput = append(YamlOutput, output)
		} else if key == Exist {
			output := ExistPath(ReplaceString(job[Exist].(string)))
			YamlOutput = append(YamlOutput, cuppago.String(output))
		} else if key == Move {
			conf := job[Move].(map[string]interface{})
			output := MovePath(ReplaceString(conf["from"].(string)), ReplaceString(conf["to"].(string)))
			YamlOutput = append(YamlOutput, output)
		} else if key == KillPortConst {
			ports := job[KillPortConst]
			if cuppago.Type(ports) == "int" {
				ports = cuppago.String(ports)
			}
			output := KillPorts(ReplaceString(ports.(string)))
			YamlOutput = append(YamlOutput, output)
		} else if key == Echo {
			output := YamlEcho(cuppago.String(job[Echo]))
			YamlOutput = append(YamlOutput, output)
		} else if key == CMD {
			output := YamlCommand(job[CMD])
			YamlOutput = append(YamlOutput, output)
		} else if key == If {
			YamlIf(job[If].(map[string]interface{}))
		} else if key == Loop {
			go YamlLoop(job[Loop].(map[string]interface{}))
		} else if key == Stop {
			os.Exit(0)
		} else if job[key] == nil {
			Log("No jobs for [", key, "]")
		} else {
			jobs := job[key].([]interface{})
			for i := 0; i < len(jobs); i++ {
				YamlJob(jobs[i].(map[string]interface{}))
			}
		}
	}
}

func YamlCommand(command interface{}) string {
	if cuppago.Type(command) == "string" {
		output := BashCommand(command.(string))
		return output
	} else {
		cmd := command.(map[string]interface{})
		dir := "./"
		if cmd["workingDirectory"] != nil {
			dir = cmd["workingDirectory"].(string)
		}
		separator := " "
		if cmd["separator"] != nil {
			separator = cmd["separator"].(string)
		}
		args := strings.Split(cmd["args"].(string), separator)
		var output string
		if cmd["background"] != nil {
			go Command(cmd["app"].(string), args, dir, "")
			output = "BACKGROUND"
		} else {
			output = Command(cmd["app"].(string), args, dir, "")
		}
		YamlOutput = append(YamlOutput, output)
		return output
	}
}

func YamlIf(data map[string]interface{}) {
	jobs := data["jobs"].([]interface{})
	if jobs == nil {
		return
	}
	output := YamlOutput[len(YamlOutput)-1]
	if data["type"] == Equal && fmt.Sprint(output) == fmt.Sprint(data["value"]) {
		YamlProcessJobs(jobs)
	} else if data["type"] == NotEqual && fmt.Sprint(output) != fmt.Sprint(data["value"]) {
		YamlProcessJobs(jobs)
	} else if data["type"] == Contain && strings.Contains(fmt.Sprint(output), fmt.Sprint(data["value"])) {
		YamlProcessJobs(jobs)
	} else if data["type"] == NotContain && !strings.Contains(fmt.Sprint(output), fmt.Sprint(data["value"])) {
		YamlProcessJobs(jobs)
	}
}

func YamlLoop(data map[string]interface{}) {
	jobs := data["jobs"].([]interface{})
	if jobs == nil {
		return
	}
	sleepTime, err := time.ParseDuration(fmt.Sprint(data["sleepTime"]) + "ms")
	if err != nil {
		cuppago.Error(err)
		return
	}
	time.Sleep(sleepTime)
	YamlProcessJobs(jobs)
	YamlLoop(data)
}

func YamlEcho(value string) string {
	value = ReplaceString(value)
	cuppago.LogFile(value)
	return value
}

func YamlCreateFile(data map[string]interface{}) string {
	file := strings.TrimSpace(ReplaceString(data["file"].(string)))
	content := ""
	if data["content"] != nil {
		content = strings.TrimSpace(ReplaceString(data["content"].(string)))
	}
	CreateFile(file, content)
	return file
}
