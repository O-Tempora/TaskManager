package main

import (
	"dip/internal/apiserver"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "./configs/apiserver.yaml", "path to config file")
}

func main() {
	flag.Parse()
	if strings.Contains(os.Args[0], "__debug_bin") {
		configPath = strings.Replace(configPath, "./", "../../", 1)
	}
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	config := apiserver.NewConfig()
	err = yaml.Unmarshal(file, config)
	if err != nil {
		log.Fatal(err)
	}

	server := apiserver.New(config)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
