package controller

import (
	"bufio"
	"fmt"
	"github.com/cloudbit-interactive/cuppago"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"strings"
	"time"
)

var YamlData map[string]interface{}
var YamlVars map[string]interface{}
var YamlOutput []string

func ProcessYamlString(yamlString string) {
	YamlOutput = []string{}
	yamlData := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(yamlString), yamlData)
	if err != nil {
		cuppago.Error(err)
	}
	YamlData = yamlData
	YamlVars = ReplaceVariables(yamlData["vars"].(map[string]interface{}))
	jobs := YamlData["jobs"].([]interface{})
	if jobs == nil {
		Log("No jobs founds")
		return
	}
	if yamlData["port"] != nil {
		port := cuppago.String(yamlData["port"])
		cuppago.Log("Continuous running in http://localhost:" + port)
		YamlProcessJobs(jobs)
		http.Handle("/favicon.ico", http.NotFoundHandler())
		http.ListenAndServe(":"+port, nil)
	} else {
		cuppago.LogFile("Process running, press [Enter] to exit...")
		YamlProcessJobs(jobs)
		bufio.NewReader(os.Stdin).ReadBytes('\n')
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
		if key == Echo {
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
		argsSeparator := " "
		if cmd["argsSeparator"] != nil {
			argsSeparator = cmd["argsSeparator"].(string)
		}
		args := strings.Split(cmd["args"].(string), argsSeparator)
		output := Command(cmd["app"].(string), args, dir)
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
