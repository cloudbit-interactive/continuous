package main

import (
	"app/src/controller"
	_ "app/src/router"
	"flag"
	"github.com/cloudbit-interactive/cuppago"
	"log"
)

func main() {
	log.SetFlags(0)
	filePath := flag.String("f", "", "Path to yaml file")
	varsPath := flag.String("v", "", "Path to yaml variables file")
	flag.Parse()
	if *filePath == "" {
		*filePath = cuppago.GetRootPath() + "/seed.yaml"
	}
	text := controller.GetFileContent(*filePath)
	controller.ProcessYamlString(text, *varsPath)
}
