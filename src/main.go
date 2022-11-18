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
	file := flag.String("file", "", " Path to yaml file")
	flag.Parse()
	if *file == "" {
		*file = cuppago.GetRootPath() + "/continuous.yaml"
	}
	text := controller.GetFileContent(*file)
	controller.ProcessYamlString(text)
}
