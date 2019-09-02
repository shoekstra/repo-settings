# repo-settings <!-- omit in toc -->

A CLI to configure repository settings across various repository hosters.

Currently supported hosters:

* GitLab

## Contents <!-- omit in toc -->

- [Why?](#why)
- [How does it work?](#how-does-it-work)
- [Installation](#installation)
- [Configuration](#configuration)
  - [GitLab](#gitlab)
    - [Project settings](#project-settings)
    - [Project integrations](#project-integrations)
      - [Slack](#slack)
- [Usage](#usage)
- [Docker](#docker)
- [Features](#features)
  - [GitLab](#gitlab-1)
- [Roadmap](#roadmap)
  - [GitLab](#gitlab-2)
- [License & Authors](#license--authors)

## Why?

When using GitLab or GitHub there are many project/repository settings you cannot set at group or organisation level.

This aims to solve that problem by letting you specify project/repository configurations in a config file and target a GitLab group or GitHub organisation and apply the config to each project/repository found.

*Note*: if a configuration can be found as a group or organisation option, it will not be supported here. This ensures the scope of this tool is only to apply settings to projects or repositories that cannot be done at group or organisation level.

## How does it work?

`repo-settings` takes settings in a config file and applies them to projects/repositorys found under a group or organisation.

It runs a few phases per project/repository found:

1. Project configuration or integration settings are stored in an object
2. This object is copied and project configuration or integration settings from the config file are applied on top of the stored object
3. If the original object and copied object differ, it is updated with the copied object

This means only settings in the config file are applied.

To give an example of this, if you have GitLab project with the Slack integration enabled and you have the Pipeline, Push and Tags events selected, and in your config file you only specify `pipeline` and `tags`, the Push event will not be disabled. `repo-settings` only ensures that what you have specified is configured, it does not ensure the state of a project/repository's settings or it's integration settings.

If you are looking for a way to configure all settings across your projects/repositories, and to ensure the state of these, you are better off using a [Terraform](https://www.terraform.io/) module to apply these configurations.

## Installation

At some point pre-built binaries will be available. Until then you will need to build it locally.

This project uses Go modules, so a minimum version of Go 1.11 is required to build the binary.

1. Clone this repository
2. Build it: `go mod download && go mod verify && go build -o repo-settings`

## Configuration

Create a (JSON or YAML) config file:

```yaml
gitlab:
  groups:
    - name: MyGroup
      integrations:
        slack:
          active: true
          events:
            - merge_request
            - pipeline
          properties:
            webhook: https://hooks.slack.com/services/T04...
            username: GitLab
```

In the case of GitLab you can also use nested groups, e.g. `MyGroup/MyNestedGroup`.

### GitLab

This section details how to configure GitLab repository settings.

#### Project settings

Coming soon!

#### Project integrations

The project integration settings are split into two parts, the Service section and the Integration Properties section.

In the example below, everything defined under the `slack` key are generic Service settings, whilst everything under the `properties` key are specific to the Slack integration:

```YAML
slack:
  active: true
  events:
    - merge_request
    - pipeline
  properties:
    webhook: https://hooks.slack.com/services/T04...
    username: GitLab
```

All integrations contain an `active` option, which should be set to `true` to enable the integration. As this tool does not delete/remove configurations it doesn't manage, an easy way to disable an integration is set `active` to `false`.

In most cases what changes between integration is which triggers are supported; each integration will list it's supported trigger events. That said. there are some integrations which do not support events and can only be configured by it's properties, as will be shown in their examples.

Generic Service settings for integrations:

| key    | description                                                     | possible settings                                         |
| ------ | --------------------------------------------------------------- | --------------------------------------------------------- |
| active | Determines if an integration is enabled                         | `true`, `false`                                           |
| events | Events that will trigger the integration, specified as an array | `["issues", "merge_request", "pipeline", "push", "tags"]` |

##### Slack

Example:

```YAML
slack:
  active: true
  events:
    - merge_request
    - pipeline
  properties:
    webhook: https://hooks.slack.com/services/T04...
    username: GitLab
```

## Usage

Specify your GitLab credentials by either exporting `GITLAB_TOKEN` and `GITLAB_URL` or using the `--gitlab-token` or `--gitlab-url` flags.

Do a dry run:

```bash
repo-settings --config config.yaml -d
```

Apply your settings:

```bash
repo-settings --config config.yaml
```

## Docker

If familiar with Docker you can use the `shoekstra/repo-settings` image, assuming you already have your variables exported locally:

```bash
docker run -it --rm -e GITLAB_TOKEN=${GITLAB_TOKEN} -e GITLAB_URL=${GITLAB_URL} -v $(pwd):/app shoekstra/repo-settings:latest -c /app/config.yaml
```

## Features

### GitLab

The following GitLab project capabilities are able to be configured:

* integrations:
  * Slack

## Roadmap

Below you find planned features, as they're completed they'll move to the section above.

Have a configuration or feature you'd like to see supported? Create a pull request that adds it to the list below.

### GitLab

* Branch protection
* Merge request approvals
* Merge request settings

## License & Authors

- Author: Stephen Hoekstra

```text
Copyright 2019 Stephen Hoekstra <stephenhoekstra@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
