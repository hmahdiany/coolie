# Pull and push container images to private registry

Coolie is a simple tool to download container images from public registries and push them in private registries.

## How it works
Coolie binary uses a config file in `yaml` format and an `env` file to load registries credentials. It pulls all image tags in config file, creates a new tag based on private registry address and finally pushes them to private registry. Coolie binary gets config file path and and environment variables file which consists of private registry information including `registry address`, `username` and `password`.

## How to use
Here is the list of environment variables used by Coolie binary:

| name | description |
|:---|:---|
| COOLIE_CONFIG | path to where config file is located |
| COOLIE_ENV | path to environment variable file |

After setting all above variables just simply execute Coolie binary file.

## Sample configuration
[Here](./configs/config.yaml) is a sample Coolie config file in `yaml` format.

[Here](./configs/env) is a sample environment variable file. Every variable in this file should starts with repository name that has been defined in `config` file. Consider following configuration:
```
repos:
  - name: quay
    source: quay.io
    destination: registry-1.example.com/repository/quay
    images:
    - name: prometheus/prometheus
      tags:
      - v2.37.9
```
With this configuration correct values in `env` file for destination registry would be like this:
```
quay_repo_address="https://registry-1.example.com"
quay_repo_username=admin
quay_repo_password=123456aA
```

Export two environment variables:

```
export COOLIE_CONFIG=configs/config.yaml
export COOLIE_ENV=configs/env
```

