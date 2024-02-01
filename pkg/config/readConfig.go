package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Repos []Repository `json:"repos"`
}

type Repository struct {
	Name        string     `json:"name"`
	Source      string     `json:"source"`
	Images      []ImageTag `json:"images"`
	Destinations []ImageDst     `json:"destinations"`
}

type ImageTag struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type ImageDst struct {
	Name string   `json:"name"`
	Address string `json:"address"`
}

func ReadConfig(path string) (cfg Config){
	var c Config

	// read config file
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)

		os.Exit(1)
	}

	// unmarshal yaml file
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
