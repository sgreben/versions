# versions - command-line version operations

`versions` is a versioning tool for the command line.

Supported operations (today):

- Compare versions
- Sort versions
- Select versions with constraints
- Fetch versions from Git tags
- Fetch versions from Docker image tags

## Contents

<!-- TOC -->

- [Contents](#contents)
- [Get it](#get-it)
- [Use it](#use-it)
    - [Compare versions](#compare-versions)
        - [Output a single boolean indicating whether one version is later than another](#output-a-single-boolean-indicating-whether-one-version-is-later-than-another)
    - [Sort versions](#sort-versions)
        - [Print versions in oldest-to-newest order](#print-versions-in-oldest-to-newest-order)
        - [Print the latest N versions in oldest-to-newest order](#print-the-latest-n-versions-in-oldest-to-newest-order)
    - [Select versions](#select-versions)
    - [Select the single latest version satisfying the given constraint](#select-the-single-latest-version-satisfying-the-given-constraint)
    - [Fetch versions](#fetch-versions)
        - [Fetch and interpret all SemVer git tags as versions](#fetch-and-interpret-all-semver-git-tags-as-versions)
        - [Fetch and determine the latest version from Git tags](#fetch-and-determine-the-latest-version-from-git-tags)
        - [Fetch and interpret all Docker image tags as versions](#fetch-and-interpret-all-docker-image-tags-as-versions)
        - [Fetch and determine the latest version from Docker image tags](#fetch-and-determine-the-latest-version-from-docker-image-tags)
    - [JSON output](#json-output)
    - [Sort order](#sort-order)
- [Licensing](#licensing)
- [Comments](#comments)

<!-- /TOC -->

## Get it

Using go get:

```sh
go get -u github.com/sgreben/versions/cmd/versions
```

Or [download the binary](https://github.com/sgreben/versions/releases/latest) from the releases page.

```sh
# Linux
curl -L https://github.com/sgreben/versions/releases/download/${VERSION}/versions_${VERSION}_linux_x86_64.tar.gz | tar xz

# OS X
curl -L https://github.com/sgreben/versions/releases/download/${VERSION}/versions_${VERSION}_osx_x86_64.tar.gz | tar xz

# Windows
curl -LO https://github.com/sgreben/versions/releases/download/${VERSION}/versions_${VERSION}_windows_x86_64.zip
unzip versions_${VERSION}_windows_x86_64.zip
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

### Select the single latest version satisfying the given constraint

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

### Fetch versions

#### Fetch and interpret all SemVer git tags as versions

```sh
$ versions --indent 2 fetch git https://github.com/sgreben/jp
```
```json
[
  {
    "Version": "1.0.0",
    "Source": {
      "Git": {
        "RepositoryURL": "https://github.com/sgreben/jp",
        "Reference": "refs/tags/1.0.0"
      }
    }
  },
  {
    "Version": "1.0.1",
    "Source": {
      "Git": {
        "RepositoryURL": "https://github.com/sgreben/jp",
        "Reference": "refs/tags/1.0.1"
      }
    }
  },
  // ...
]
```

#### Fetch and determine the latest version from Git tags

```sh
$ versions fetch -l 1 git https://github.com/sgreben/jp
```
```json
[{"Version":"1.1.11","Source":{"Git":{"RepositoryURL":"https://github.com/sgreben/jp","Reference":"refs/tags/1.1.11"}}}]
```


#### Fetch and interpret all Docker image tags as versions

```sh
$ versions fetch docker alpine | jq .
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
