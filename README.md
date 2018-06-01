# versions

`versions` is a tool for working with (SemVer) versions.

Supported operations (today):

- Compare versions
- Sort versions
- Select versions with constraints
- Fetch versions from Git tags
- Fetch versions from Docker image tags

Planned operations for `v1.0.0`:

- Dependency version selection using [MVS](https://research.swtch.com/vgo-mvs)
- Dependency version graph analysis (without solving)

## Contents

- [Get it](#get-it)
- [Use it](#use-it)
    - [Compare versions](#compare-versions)
    - [Sort versions](#sort-versions)
    - [Select versions](#select-versions)
    - [Fetch versions](#fetch-versions)
    - [JSON output](#json-output)
- [Licensing](#licensing)
- [Comments](#comments)

## Get it

Using go get:

```sh
go get -u github.com/sgreben/versions/cmd/versions
```

Or [download the binary](https://github.com/sgreben/versions/releases/latest) from the releases page.

```sh
# Linux
curl -L https://github.com/sgreben/versions/releases/download/0.0.1/versions_0.0.1_linux_x86_64.tar.gz | tar xz

# OS X
curl -L https://github.com/sgreben/versions/releases/download/0.0.1/versions_0.0.1_osx_x86_64.tar.gz | tar xz

# Windows
curl -LO https://github.com/sgreben/versions/releases/download/0.0.1/versions_0.0.1_windows_x86_64.zip
unzip versions_0.0.1_windows_x86_64.zip
```

Also available as a [docker image](https://quay.io/repository/sergey_grebenshchikov/versions?tab=tags):

```sh
docker run quay.io/sergey_grebenshchikov/versions
```

## Use it

```text
Usage: versions COMMAND [arg...]

do things with versions

Options:
      --indent   Set the indentation of JSON output (default 0)
  -q, --quiet    Disable all log output (stderr)
  -s, --silent   Disable all log output (stderr) and all normal output (stdout)

Commands:
  sort         Sort versions
  compare      Compare versions
  fetch        Fetch versions
  select       Select versions given constraints
  complete     Shell completion (zsh, fish, bash)
  help         Display help for a command

Run 'versions COMMAND --help' for more information on a command.
```

### Compare versions

```text
Usage: versions compare [OPTIONS] COMMAND [arg...]

Compare versions

Options:
      --fail   Exit with non-zero code if the result is 'false'

Commands:
  later        Check if a version is strictly later than another version
  earlier      Check if a version is strictly earlier than another version

Run 'versions compare COMMAND --help' for more information on a command.
```

#### Output a single boolean indicating whether one version is later than another

```sh
$ versions compare later 1.0.0 0.1.0
```
```json
true
```

```sh
$ versions compare later 1.0.0 2.1.0
```
```json
false
```

### Sort versions

```text
Usage: versions sort [OPTIONS] [VERSIONS...]

Sort versions

Arguments:
  VERSIONS       Versions to sort

Options:
  -l, --latest   Print only the latest `N` versions (default 0)
```

#### Print versions in oldest-to-newest order

```sh
$ versions sort 2.0.0 0.1.0 10.0.0
```
```json
["0.1.0","2.0.0","10.0.0"]
```

#### Print the latest N versions in oldest-to-newest order

```sh
$ versions --latest=2 sort 2.0.0 0.1.0 10.0.0
```
```json
["2.0.0","10.0.0"]
```

### Select versions

```text
Usage: versions select [OPTIONS] COMMAND [arg...]

Select versions given constraints

Options:
      --from-git      Fetch candidate versions from Git tags
      --from-docker   Fetch candidate versions from Docker tags

Commands:
  single              Select a single version
  graph               Select versions to satisfy a constraint graph

Run 'versions select COMMAND --help' for more information on a command.
```

#### Select the single latest version satisfying the given constraint

```sh
$ versions select single '2.*.*' 2.0.0 0.1.0 10.0.0
```
```json
"2.0.0"
```

```sh
$ versions select single '*' 2.0.0 0.1.0 10.0.0
```
```json
"10.0.0"
```

```sh
$ versions select single '^0.0.1' 2.0.0 0.1.0 10.0.0
```
```json
"0.1.0"
```

#### Select the single latest version from Git tags satisfying the given constraint

```sh
$ versions select --from-git=https://github.com/sgreben/jp single '~1.0.0'
```
```json
'1.0.1"
```

```sh
$ versions select --from-git=https://github.com/sgreben/jp single '^1.0.0'
```
```json
"1.1.11"
```

#### Select the single latest version from Docker tags satisfying the given constraint

```sh
$ versions select --from-docker=alpine single '<3.7'
```
```json
"3.6.0"
```

```sh
$ versions select --from-docker=alpine single '^3.0.0'
```
```json
"3.7.0"
```


### Fetch versions

```text
Usage: versions fetch [OPTIONS] COMMAND [arg...]

Fetch versions

Options:
  -l, --latest   Print only the latest `N` versions (default 0)

Commands:
  git            Fetch versions from Git tags
  docker         Fetch versions from Docker image tags

Run 'versions fetch COMMAND --help' for more information on a command.
```

#### Fetch and interpret all SemVer git tags as versions

```sh
$ versions --indent=2 fetch git https://github.com/sgreben/jp
```
```json
[
  {
    "Version": "1.0.0",
    "Source": {
      "Git": {
        "Repository": {
          "URL": "https://github.com/sgreben/jp"
        },
        "Reference": "refs/tags/1.0.0"
      }
    }
  },
  {
    "Version": "1.0.1",
    "Source": {
      "Git": {
        "Repository": {
          "URL": "https://github.com/sgreben/jp"
        },
        "Reference": "refs/tags/1.0.1"
  // ...
]
```

#### Fetch and determine the latest version from Git tags

```sh
$ versions fetch -l 1 git https://github.com/sgreben/jp
```
```json
[{"Version":"1.1.11","Source":{"Git":{"Repository":{"URL":"https://github.com/sgreben/jp"},"Reference":"refs/tags/1.1.11"}}}]
```


#### Fetch and interpret all Docker image tags as versions

```sh
$ versions --indent=2 fetch docker alpine
```
```json
[
  {
    "Version": "2.6.0",
    "Source": {
      "Docker": {
        "Image": "library/alpine:2.6",
        "Tag": "2.6"
      }
    }
  },
  {
    "Version": "2.7.0",
    "Source": {
      "Docker": {
        "Image": "library/alpine:2.7",
        "Tag": "2.7"
      }
    }
  },
  // ...
]
```

#### Fetch and determine the latest version from Docker image tags

```sh
$ versions fetch -l 1 docker alpine
```
```json
[{"Version":"3.7.0","Source":{"Docker":{"Image":"library/alpine:3.7","Tag":"3.7"}}}]
```

### JSON output

The default output format is JSON, one value per line:

```sh
$ versions sort 0.10 0.2 1.0 1.1 1.1.1-rc1 1.1.1
```
```json
["0.2.0","0.10.0","1.0.0","1.1.0","1.1.1-rc1","1.1.1"]
```

To output multi-line indented JSON, specify a value for the `--indent` option:


```sh
$ versions --indent=2 sort 0.10 0.2 1.0 1.1 1.1.1-rc1 1.1.1
```
```json
[
  "0.2.0",
  "0.10.0",
  "1.0.0",
  "1.1.0",
  "1.1.1-rc1",
  "1.1.1"
]
```

### Sort order

All commands that produce sorted lists of versions produce them in the **oldest-first**, **latest-last** order:

```sh
$ versions sort 0.0.1 1.0.0
```
```json
["0.0.1","1.0.0"]
```

## Licensing

- Any original code is licensed under the [MIT License](./LICENSE).
- The included version of [github.com/Masterminds/semver](https://github.com/Masterminds/semver) is licensed under [what looks like the MIT license](https://github.com/Masterminds/semver/blob/c7af12943936e8c39859482e61f0574c2fd7fc75/LICENSE.txt).
- Included portions of [github.com/kubernetes/client-go](https://github.com/kubernetes/client-go/tree/master/util/jsonpath) are licensed under the Apache License 2.0.

## Comments

Feel free to [leave a comment](https://github.com/sgreben/versions/issues/1) or create an issue.
