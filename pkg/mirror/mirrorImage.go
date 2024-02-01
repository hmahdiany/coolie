package mirror

import (
	"bufio"
	"context"

	//"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/hmahdiany/coolie/pkg/config"
	"github.com/pkg/errors"
)

type outputLog struct {
	Status string `json:"status"`
	Id     string `json:"id"`
}

func MirrorImages(ctx context.Context, cfg config.Config, imageMap map[string][]string) []error {
	errorList := []error{}

	dockerClient, err := client.NewClientWithOpts()
	if err != nil {
		errorList = append(errorList, errors.Wrap(err, "failed to create docker client"))
		return errorList
	}

	//login to destination registries
	logintoregistry(ctx, cfg, dockerClient)

	wg := sync.WaitGroup{}

	for key, value := range imageMap {
		wg.Add(1)
		go func(cfg config.Config, source string, destination []string) {
			defer wg.Done()

			err := mirrorImage(ctx, cfg, dockerClient, source, destination)

			if err != nil {
				printingError := errors.Wrapf(err, "failed to mirror image %s to %s", source, destination)
				fmt.Println(printingError)
			}

		}(cfg, key, value[:])
		wg.Wait()
	}
	return errorList
}

func mirrorImage(ctx context.Context, cfg config.Config, dockerClient *client.Client, source string, destination []string) error {
	if err := pullImage(ctx, dockerClient, source); err != nil {
		return errors.Wrapf(err, "failed to pull image %s", source)
	}
	if err := renameImage(ctx, dockerClient, source, destination); err != nil {
		return errors.Wrapf(err, "failed to rename image %s to %s", source, destination)
	}
	if err := pushImage(ctx, cfg, dockerClient, destination); err != nil {
		return errors.Wrapf(err, "failed to push image %s", destination)
	}
	return nil
}

func pullImage(ctx context.Context, dockerClient *client.Client, imageTag string) error {

	fmt.Printf("start pulling image %s\n", imageTag)

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

func renameImage(ctx context.Context, dockerClient *client.Client, source string, destination []string) error {
	var err error
	for _, dst := range destination {
		err = dockerClient.ImageTag(ctx, source, dst)
	}
	return err
}

func pushImage(ctx context.Context, cfg config.Config, dockerClient *client.Client, imageTags []string) error {


	// registryAuthConfigBytes, err := json.Marshal(registryAuthConfig)
	// if err != nil {
	// 	return errors.Wrap(err, "failed to create regstry auth config")
	// }

	// registryAuthConfigEncoded := base64.URLEncoding.EncodeToString(registryAuthConfigBytes)

	// opts := types.ImagePushOptions{
	// 	RegistryAuth: registryAuthConfigEncoded,
	// }

	// loop through destinationRegistryList to push images
	for _, tag := range imageTags {

		fmt.Println("tag list: ", imageTags)
		fmt.Println("current tag: ", tag)
		fmt.Printf("Starting push loop: %v\n", tag)
		rc, err := dockerClient.ImagePush(ctx, tag, types.ImagePushOptions{})
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

		fmt.Printf("pushed image %s\n", tag)
	}

	return nil
}

func CreateImageMap(cfg config.Config) map[string][]string {
	imageNames := map[string][]string{}

	// create source and destination image name
	for _, registry := range cfg.Repos {
		for _, image := range registry.Images {
			for _, tag := range image.Tags {
				sourceImageTag := registry.Source + "/" + image.Name + ":" + tag
				for _, dst := range registry.Destinations {
					desinationImageTag := dst.Address + "/" + image.Name + ":" + tag
					imageNames[sourceImageTag] = append(imageNames[sourceImageTag], desinationImageTag)
				}
			}
		}
	}

	return imageNames
}
