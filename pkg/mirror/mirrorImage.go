package mirror

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/hmahdiany/coolie/pkg/config"
	"github.com/pkg/errors"
)

type outputLog struct {
	Status string `json:"status"`
	Id     string `json:"id"`
}

func MirrorImages(ctx context.Context, imageMap map[string]string) []error {
	errorList := []error{}

	dockerClient, err := client.NewClientWithOpts()
	if err != nil {
		errorList = append(errorList, errors.Wrap(err, "failed to create docker client"))
		return errorList
	}

	wg := sync.WaitGroup{}

	for key, value := range imageMap {
		wg.Add(1)
		go func(source, destination string) {
			defer wg.Done()

			err := mirrorImage(ctx, dockerClient, source, destination)

			if err != nil {
				printingError := errors.Wrapf(err, "failed to mirror image %s to %s", source, destination)
				fmt.Println(printingError)
			}

		}(key, value)
		wg.Wait()
	}
	return errorList
}

func mirrorImage(ctx context.Context, dockerClient *client.Client, source, destination string) error {
	if err := pullImage(ctx, dockerClient, source); err != nil {
		return errors.Wrapf(err, "failed to pull image %s", source)
	}
	if err := renameImage(ctx, dockerClient, source, destination); err != nil {
		return errors.Wrapf(err, "failed to rename image %s to %s", source, destination)
	}
	if err := pushImage(ctx, dockerClient, destination); err != nil {
		return errors.Wrapf(err, "failed to push image %s", destination)
	}
	return nil
}

func pullImage(ctx context.Context, dockerClient *client.Client, imageTag string) error {

	rc, err := dockerClient.ImagePull(ctx, imageTag, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	var outputlog outputLog
	scanner := bufio.NewScanner(rc)

	for scanner.Scan() {
		line := scanner.Text()
		err := json.Unmarshal([]byte(line), &outputlog)
		if err != nil {
			return err
		}

		if outputlog.Status == "Download complete" {
			fmt.Printf("status: %v, id: %v\n", outputlog.Status, outputlog.Id)
		}
	}

	fmt.Printf("pulled image %s\n", imageTag)
	return nil
}

func renameImage(ctx context.Context, dockerClient *client.Client, source, destination string) error {
	return dockerClient.ImageTag(ctx, source, destination)
}

func pushImage(ctx context.Context, dockerClient *client.Client, imageTag string) error {

	var registryAuthConfig = registry.AuthConfig{
		Username:      os.Getenv("REPO_USERNAME"),
		Password:      os.Getenv("REPO_PASSWORD"),
		ServerAddress: os.Getenv("REPO_ADDRESS"),
	}

	registryAuthConfigBytes, err := json.Marshal(registryAuthConfig)
	if err != nil {
		return errors.Wrap(err, "failed to create regstry auth config")
	}
	registryAuthConfigEncoded := base64.URLEncoding.EncodeToString(registryAuthConfigBytes)

	opts := types.ImagePushOptions{
		RegistryAuth: registryAuthConfigEncoded,
	}

	rc, err := dockerClient.ImagePush(ctx, imageTag, opts)
	if err != nil {
		return errors.Wrap(err, "failed to push image with docker client")
	}

	var outputlog outputLog
	scanner := bufio.NewScanner(rc)

	for scanner.Scan() {
		line := scanner.Text()
		err := json.Unmarshal([]byte(line), &outputlog)
		if err != nil {
			return err
		}

		fmt.Printf("status: %v, id: %v\n", outputlog.Status, outputlog.Id)
	}
	
	fmt.Printf("pushed image %s\n", imageTag)
	return nil
}

func CreateImageMap(cfg config.Config) map[string]string {
	imageNames := map[string]string{}

	// create source and destination image name
	for _, registry := range cfg.Repos {
		for _, image := range registry.Images {
			for _, tag := range image.Tags {
				sourceImageTag := registry.Source + "/" + image.Name + ":" + tag
				desinationImageTag := registry.Destination + "/" + image.Name + ":" + tag

				imageNames[sourceImageTag] = desinationImageTag
			}
		}
	}

	return imageNames
}
