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