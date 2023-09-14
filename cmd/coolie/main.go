package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/hmahdiany/coolie/pkg/config"
	"github.com/hmahdiany/coolie/pkg/mirror"
)

func main() {

	//Get config file path from env
	path := os.Getenv("COOLIE_CONFIG")
	if len(path) == 0 {
		fmt.Println("Config file is not specified. Set COOLIE_CONFIG environment variable")
		os.Exit(1)
	} else {
		_, err := os.Stat(path)
		if err != nil {
			log.Printf("Config file %v does not exist\n.", path)
			os.Exit(2)
		}

		fmt.Printf("Using %v as config file\n", path)
	}

	// get config value
	cfg := config.ReadConfig(path)

	// create a map from container images in config file
	imageMap := mirror.CreateImageMap(cfg)

	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancelFunc()

	// start mirroring
	errorList := mirror.MirrorImages(ctx, imageMap)

	if len(errorList) == 0 {
		fmt.Println("All image tags in config file were processed")
		return
	}

	for _, err := range errorList {
		fmt.Println(err)

	}
}
