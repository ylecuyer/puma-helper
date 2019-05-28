# How to release puma-helper

The release process is done using [Goreleaser](https://goreleaser.com/).

Goreleaser configuration file is already created and formatted. it can be found as `.goreleaser.yml` at the root path of the repository.

## Prerequisites starter pack

* Clone the project within your `GOPATH` under `src/github.com/dimelo/puma-helper`.
* Follow [the official install guide](https://goreleaser.com/install/) based on the OS you want to build the binary.
* Get and export a `GITHUB_TOKEN` [just here](https://github.com/settings/tokens), which should contain a valid GitHub token with the repo scope, in your env bash/zsh.
* Install `rpm` package from your package manager.
* Go get `github.com/goreleaser/nfpm` and install using `go install` command.

## Usage guide

When the prerequisites install is done, you're ready to go and release `puma-helper`.

### Tag your commits

In first, tag your commits as a release using [semantic versioning](https://semver.org/) and push it.
```bash
$ git tag -a v0.1.0 -m "First release"
$ git push origin v0.1.0
````

### Dry release

Your new tag is now push on Github, you can [found it here](https://github.com/dimelo/puma-helper/tags) but without any context or binaries. This is the next step.

Just follow the next lines to check if everything works fine, before doing the final release process.

```bash
goreleaser release --skip-publish --rm-dist
```

`--rm-dist` option is here to clean up your build folder.

### Final release

When the build release result is `succeeded`, you just have to follow in the same way the last command without the `--skip-publish` option.

```bash
goreleaser release --rm-dist
```

If succeeded, the new release [will appear here](https://github.com/dimelo/puma-helper/releases). 

### Extra bonus

Feel free to surname your release just [like this](https://github.com/dimelo/puma-helper/releases/tag/v1.0.0) and add a cool pic of it related to the name :rocket:
