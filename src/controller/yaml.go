package controller

import (
	"github.com/cloudbit-interactive/cuppago"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
)

var YamlData map[string]interface{}
var YamlVars map[string]interface{}

func ProcessYamlString(yamlString string) {
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
	cuppago.Log(job)
	if job[CMD] != nil {
		cmd := cuppago.String(job[CMD])
		cmd = YamlReplaceString(cmd)
		cuppago.Log("Process COMMAND", cmd)

		//output := exec.Command("if", "[ -f ${serviceDir}/continuous.service ] ; then echo 1 ; else echo 0 ; fi")
		//cuppago.Log("OUTPUT", output)

	} else if job[IF] != nil {
		cuppago.Log("Process IF", job[IF])
	} else if job[STOP] != nil {
		cuppago.Log("Process STOP", job[STOP])
		os.Exit(0)
	}
}

func YamlReplaceString(string string) string {
	string = cuppago.ReplaceNotCase(string, "{serviceDir}", "ddd")
	return string
}
