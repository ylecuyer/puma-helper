# puma-helper

Puma-helper aims to implement missing centralized and human readeable features from puma unix socket in one place.

## Disclaimer

! Work and documentation still in progress !

## Installation

Download and install the latest [release](https://github.com/dimelo/puma-helper/releases) packaged for your distribution.

## Configuration

You can find a `.puma-helper.yaml` example in this repository or use `init` command directly.

The configuration file can be read by the binary from the install directory or `$HOME` path, if defined.

Only the `application name` and `path` are mandatory.

```yaml
# .puma-helper.yaml
applications:
  # Mandatory - string
  # Your application name
  "your_app_name":
    # Mandatory - string
    # Path to your application
    path : "/home/path/to/your/app"

    # Mandtory - array of string
    # Path(s) to puma state file(s)
    pumastatepaths:
    - /home/your_app/current/tmp/pids/puma.state
    - /home/your_app/current/tmp/pids/puma_api.state

    # Optional - string
    # Description or informations related to the application
    # Default description: ""
    description : "Related to my super app, in production"

    # Optional - int
    # Active thread warn and critical % must be > 1 and < 100
    # Default Active thread warn: 50
    # Default Active thread critical: 80
    #thread_warn : 50
    #thread_critical : 80

    # Optional - int
    # CPU warn and critical % must be > 1 and < 100
    # Default CPU warn: 50
    # Default CPU critical: 80
    #cpu_warn : 50
    #cpu_critical : 80

    # Optional - int
    # Memory warn and critical usage must be > 0
    # Default memory warn: 500
    # Default memory critical: 1000
    #memory_warn : 500
    #memory_critical : 1000
```

## CLI usage

### Init

Init command permit to create configuration file if it doesn't exist (or replace it).

Simply follow the questions and enter absolute path only.

At the end, the configuration file will be placed under `$HOME/.puma-helper.yaml`.

Run the command:
```
puma-helper init
```

### Status

Status command permit to centralize puma unix socket status metrics.

Run the command:
```
puma-helper status
```

#### Options

* `filter`: Only show applications who match with given string
* `json`: Return JSON object who contains all informations
* `details`: Show more details about apps and workers

## Release the CLI

Checkout [RELEASE.md](RELEASE.md)

## Report a bug

Directly open an issue and follow given steps.

## Hacking

Directly open an PR and follow given steps.

## License

This project is under the MIT license, see the [LICENSE](LICENSE) file for details.