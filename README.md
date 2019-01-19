# puma-helper

Puma-helper aims to implement missing centralized and human readeable features from pumactl in one place.

## Disclaimer

! Work and documentation still in progress !

## Installation

Download and install the latest [release](https://github.com/dimelo/puma-helper/releases) packaged for your distribution.

## Configuration

You can find a `puma-helper.yaml` example in this repository.

The configuration file can be read by the binary from the install directory or `$HOME` path, if defined.

Only the `application name` and `path` are mandatory.

```yaml
# puma-helper.yaml
applications:
  # Mandatory - string
  # Your application name
  "your_app_name":
    # Mandatory - string
    # Path to your application
    path : "/home/path/to/your/app"

    # Optional - string
    # Description or informations related to the application
    # Default description: ""
    description : "Related to my super app, in production"

    # Optional - string
    # Pumactl and Puma state path to files must be absolute
    # Default pumactlpath: current/bin/pumactl (from path of your app)
    # Default pumastatepath: /tmp/pids/puma.state (from path of your app)
    #pumactlpath : "/home/path/to/pumactl"
    #pumastatepath : "/home/path/to/puma_state"

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

### Status

Status command permit to centralize pumactl status metrics

Run the command:
```
puma-helper status
```

#### Options

* `filter`: Only show applications who match with given string

## Report a bug

Directly open an issue and follow given steps

## Hacking

Directly open an PR and follow given steps

## License

This project is under the MIT license, see the [LICENSE](LICENSE) file for details.