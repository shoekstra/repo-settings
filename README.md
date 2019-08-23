# repo-settings

A CLI to configure repository settings across various repository hosters.

Currently supported hosters:

* GitLab

## Why?

When using GitLab or GitHub there are many project/repository settings you cannot set at group or organisation level. This tool aims to solve that problem by letting you specify project/repository configurations in a config file and target a GitLab group or GitHub organisation. It will then look up all projects/repositories under that group or organisation and apply the specified config to each.

- [Why?](#why)
- [Installation](#installation)
- [Configuration](#configuration)
  - [GitLab](#gitlab)
    - [Project settings](#project-settings)
    - [Project integrations](#project-integrations)
      - [Slack](#slack)
- [Usage](#usage)
- [Features](#features)
  - [GitLab](#gitlab-1)
- [Roadmap](#roadmap)
  - [GitLab](#gitlab-2)
- [License & Authors](#license--authors)

## Installation

At some point pre-built binaries will be available. Until then you will need to build it locally.

This project uses Go modules, so a minimum version of Go 1.11 is required to build the binary.

1. Clone this repository
2. `go build -o repo-settings`

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

For example:

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

Everything defined under the `slack` key are generic Service settings, whilst everything under the `properties` key are specific to the Slack integration.

Generic Service settings that apply to all integrations:

| key    | description                                                     | valid settings                                            |
| ------ | --------------------------------------------------------------- | --------------------------------------------------------- |
| active | Determines if an integration is enabled                         | `true`, `false`                                           |
| events | Events that will trigger the integration, specified as an array | `["issues", "merge_request", "pipeline", "push", "tags"]` |

With this in mind, integration settings shown below in this section only apply to the `properties` key within the integration section.

##### Slack

Example:

```YAML
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

## Features

### GitLab

The following GitLab project capabilities are able to be configured:

* Slack integration

## Roadmap

Below you find planned features, as they're completed they'll move to the section above.

Have a configuration or feature you'd like to see supported? Create a pull request that adds it to the list below.

### GitLab

* Branch protection
* Required approvers

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
