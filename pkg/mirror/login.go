package mirror

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/hmahdiany/coolie/pkg/config"
	"github.com/pkg/errors"
)

func logintoregistry(ctx context.Context, cfg config.Config, dockerClient *client.Client) {

	// get destination registries
	destinationRegistryList := config.DestinationRegistry(cfg)

	// login to all destination registries
	for _, reg := range destinationRegistryList {
		fmt.Printf("Login to: %v\n", reg)
		var registryAuthConfig = registry.AuthConfig{
			Username:      os.Getenv(reg + "_repo_username"),
			Password:      os.Getenv(reg + "_repo_password"),
			ServerAddress: os.Getenv(reg + "_repo_address"),
		}

		_, err := dockerClient.RegistryLogin(ctx, registryAuthConfig)
		if err != nil {
			printingError := errors.Wrapf(err, "failed to login to %s", reg)
			fmt.Println(printingError)
		}
	}

}
