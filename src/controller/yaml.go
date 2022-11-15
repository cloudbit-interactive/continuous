package controller

import (
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
	port := "9323"
	if yamlData["port"] != nil {
		port = cuppago.String(yamlData["port"])
	}
	YamlData = yamlData
	YamlVars = yamlData["vars"].(map[string]interface{})
	jobs := YamlData["jobs"].([]interface{})
	if jobs != nil {
		YamlProcessJobs(jobs)
	}
	cuppago.Log("Continuous running in http://localhost:" + port)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":"+port, nil)
}

func YamlProcessJobs(jobs []interface{}) {
	for i := 0; i < len(jobs); i++ {
		YamlJob(jobs[i].(map[string]interface{}))
	}
}

func YamlJob(job map[string]interface{}) {
	for key := range job {
		if key == CMD {
			output := YamlCommand(job[CMD])
			YamlOutput = append(YamlOutput, output)
		} else if key == If {
			YamlIf(job[If].(map[string]interface{}))
		} else if key == Loop {
			go YamlLoop(job[Loop].(map[string]interface{}))
		} else if key == Stop {
			os.Exit(0)
		} else if job[key] == nil {
			cuppago.LogFile("No jobs for [", key, "]")
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
			dir = ReplaceString(dir)
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
	YamlProcessJobs(jobs)
	sleepTime, err := time.ParseDuration(fmt.Sprint(data["sleepTime"]) + "ms")
	if err != nil {
		cuppago.Error(err)
		return
	}
	time.Sleep(sleepTime)
	YamlLoop(data)
}
