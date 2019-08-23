# repo-settings

A simple CLI to configure repositories settings across various repository hosters.

It reads a configuration file containing a GitHub organisation (not yet supported) and/or GitLab group and will configure all repositories found within with the defined settings.

To get started, create a configuration file and pass the --config option.

- [Installation](#installation)
- [Usage](#usage)
  - [GitLab](#gitlab)
    - [Project settings](#project-settings)
    - [Project integrations](#project-integrations)
      - [Slack](#slack)
- [License & Authors](#license--authors)

## Installation

At some later point pre-built binaries will be available. Until then you will need to build it manually.

This project uses Go modules, so a minimum version of Go 1.11 is required to build the binary.

1. Clone this repository
2. `go build -o repo-settings`

## Usage

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

Do a dry run:

```bash
repo-settings --config config.yaml -d
```

Apply your settings:

```bash
repo-settings --config config.yaml
```

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

| key    | description                              | valid settings                                            |
| ------ | ---------------------------------------- | --------------------------------------------------------- |
| active | Determines if an integration is enabled  | `true`, `false`                                           |
| events | Events that will trigger the integration | `["issues", "merge_request", "pipeline", "push", "tags"]` |

With this in mind, integration settings shown below in this section only apply to the `properties` key within the integration section.

##### Slack

Example:

```YAML
webhook: https://hooks.slack.com/services/T04...
username: GitLab
```

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
