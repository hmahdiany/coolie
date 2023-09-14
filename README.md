# Pull and push container images to private registry

Coolie is a simple tool to download container images from public registries and push them in private registries.

## How it works
Coolie binary uses a config file in `yaml` format. It pulls all image tags in config file, creates a new tag based on private registry address and finally pushes them to private registry. Coolie binary gets config file path and private registry information including `registry address`, `username` and `password` all via environment variables.

## How to use
Here is the list of environment variables used by Coolie binary:

| name | description |
|:---|:---|
| COOLIE_CONFIG | path to where config file is located |
| REPO_USERNAME | registry username |
| REPO_PASSWORD | registry password |
| REPO_ADDRESS | registry address to login |

After setting all above variables just simply execute Coolie binary file.

## Sample configuration
[Here](./configs/config.yaml) is a sample Coolie config file in `yaml` format.

Export all four environment variables:

```
export COOLIE_CONFIG=configs/config.yaml
export REPO_USERNAME=coolie
export REPO_PASSWORD=123456aA
export REPO_ADDRESS="https://registry.example.com"
```

