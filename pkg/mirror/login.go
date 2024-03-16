package mirror

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/hmahdiany/coolie/pkg/config"
	"github.com/pkg/errors"
)

func logintoregistry(ctx context.Context, cfg config.Config, dockerClient *client.Client, repo string) (types.ImagePushOptions, error) {

	var opts types.ImagePushOptions

	// get destination registries
	destinationRegistry := config.DestinationRegistry(cfg)

	// login to all destination registries
	for k, v := range destinationRegistry {
		if repo == v {
			fmt.Printf("Login to: %v\n", k)
			var registryAuthConfig = registry.AuthConfig{
				Username:      os.Getenv(k + "_repo_username"),
				Password:      os.Getenv(k + "_repo_password"),
				ServerAddress: os.Getenv(k + "_repo_address"),
			}

			_, err := dockerClient.RegistryLogin(ctx, registryAuthConfig)
			if err != nil {
				return types.ImagePushOptions{}, errors.Wrapf(err, "failed to login to %s", k)
			}

			registryAuthConfigBytes, err := json.Marshal(registryAuthConfig)
			if err != nil {
				return types.ImagePushOptions{}, errors.Wrap(err, "failed to create regstry auth config")
			}

			registryAuthConfigEncoded := base64.URLEncoding.EncodeToString(registryAuthConfigBytes)

			opts = types.ImagePushOptions{
				RegistryAuth: registryAuthConfigEncoded,
			}
		}
	}

	return opts, nil
}
