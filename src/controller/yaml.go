package controller

import (
	"github.com/cloudbit-interactive/cuppago"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"strings"
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
	YamlProcessMainJobs()
	//YamlMainJobs(yamlData["jobs"].([]interface{}))
	cuppago.Log("Continuous running in http://localhost:" + port)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":"+port, nil)
}

func YamlProcessMainJobs() {
	if YamlData["jobs"] == nil {
		return
	}
	jobs := YamlData["jobs"].([]interface{})
	for i := 0; i < len(jobs); i++ {
		YamlJob(jobs[i].(map[string]interface{}))
	}
}

func YamlJob(job map[string]interface{}) {
	//cuppago.Log("NEW-->", job, cuppago.Type(job))
	for key := range job {
		//cuppago.Log(job)
		if job[key] == nil {
			println("no subjobs for [", key, "]")
		} else {
			jobs := job[key].([]interface{})
			for i := 0; i < len(jobs); i++ {
				YamlProcessJob(jobs[i].(map[string]interface{}))
			}
		}
	}
}

func YamlProcessJob(job map[string]interface{}) {
	if job[CMD] != nil {
		if cuppago.Type(job[CMD]) == "string" {
			output := BashCommand(job[CMD].(string))
			YamlOutput = append(YamlOutput, output)
		} else {
			cmd := job[CMD].(map[string]interface{})
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
		}
		cuppago.Log(YamlOutput)
	} else if job[IF] != nil {
		cuppago.Log("Process IF", job[IF])
	} else if job[STOP] != nil {
		cuppago.Log("Process STOP", job[STOP])
		os.Exit(0)
	}
}

func YamlCommand(string string) string {
	return ""
}
