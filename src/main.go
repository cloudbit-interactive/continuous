package main

import (
	"app/src/controller"
	_ "app/src/router"
	"flag"
	"log"
)

func main() {
	log.SetFlags(0)
	file := flag.String("file", "", " Path to yaml file")
	flag.Parse()
	if *file == "" {
		*file = "continuous.yaml"
	}
	text := controller.GetFileContent(*file)
	controller.ProcessYamlString(text)
}
