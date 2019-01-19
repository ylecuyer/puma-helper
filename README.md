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
  "your_app_name":
    path : "/home/path/to/your/app"
    description : "Related to my super app, in production"
    #pumactlpath : "/home/path/to/pumactl"
    #pumastatepath : "/home/path/to/puma_state"
    # Active thread warn and critical % must be > 1 and < 100
    # Default active thread warn: 50, critical: 80
    #thread_warn : 50
    #thread_critical : 80
    # CPU warn and critical % must be > 1 and < 100
    # Default CPU warn: 50, critical: 80
    #cpu_warn : 50
    #cpu_critical : 80
    # Memory usage warn and critical must be > 0
    # Default memory warn: 500, critical: 1000
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