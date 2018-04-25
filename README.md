# docker-tools

[![Build Status](https://travis-ci.org/pengsrc/docker-tools.svg?branch=master)](https://travis-ci.org/pengsrc/docker-tools)
[![Go Report Card](https://goreportcard.com/badge/github.com/pengsrc/docker-tools)](https://goreportcard.com/report/github.com/pengsrc/docker-tools)
[![License](http://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/pengsrc/docker-tools/blob/master/LICENSE)

Handy tools for Docker.

## What's Inside

Currently, this project contains a command-line tool called `docker-tools`.

And some useful scripts will be provided later on.

## Installation the Command Line Tool

### Install from Source Code

``` bash
$ git clone git@github.com:pengsrc/docker-tools.git
$ cd docker-tools
$ glide install
...
[INFO]	Replacing existing vendor dependencies
$ make install
...
Installing into /data/go/bin/docker-tools...
Done
```

### Download Precompiled Binary

1. Go to [releases tab](https://github.com/pengsrc/docker-tools/releases) and download the binary for your operating system.
2. Unarchive the downloaded file, and put the executable file `docker-tools` into a directory that in the `$PATH` environment variable, for example `/usr/local/bin`.
3. Run `docker-tools --help` to get started.

``` Bash
$ docker-tools --help
Handy tools for Docker

Usage:
  docker-tools [flags]
  docker-tools [command]

Available Commands:
  help          Help about any command
  remote-import Import docker image from one registry to another
  remote-build  Build docker image and push to specified registry

Flags:
  -c, --config string   Configuration file (default is ${HOME}/.docker-tools.yaml)
      --help            Show help
  -v, --version         Show version

Use "docker-tools [command] --help" for more information about a command.
```

## Usage of the Command Line Tool

### Configuration

Basically, you don't need a config file, because all of the options can be passed from command flags. But, you can also place a config file to predefine some options that are frequently used. The default config file location is `~/.docker-tools.yaml`, alternatively, the config file location can be override with `--config` flag.

A config file example can be found at `./cmds/configs/docker-tools.yaml.example`:

``` YAML
# Configuration example for docker-tools
# than should be placed at "~/.docker-tools.yaml".

# Host to remote build & import images.
builder:
  host: 127.0.0.1
  port: 22

# Docker registry to push built images and imported images.
registry: registry.example.com
```


### Build Docker Image on Remote Machine

Subcommand `remote-build` of `docker-tools` has the ability to build Docker image on remote server and push the built image to a specified image registry. It can also detect git repository and extract corresponding version of source code to build.

All available options are:

``` Bash
$ docker-tools remote-build --help
...
Flags:
  -a, --after string          Command to execute after the build
  -b, --before string         Command to execute before the build
  -d, --directory string      Source directory to use
  -f, --dockerfile string     Dockerfile to use in build (default "Dockerfile")
  -e, --exclude stringArray   Files to exclude in package
  -g, --git-archive           Use git archive to pack files
      --help                  Show help
  -h, --host string           SSH host to run import procedures (default "127.0.0.1")
  -i, --include stringArray   Files to include in package
  -p, --port int              SSH port to connect (default 22)
  -r, --registry string       Registry to push image (default "registry.example.com")
  -u, --user string           SSH username (default "root")
...
```

Example:

``` Bash
$ glide install
...
[INFO]	Replacing existing vendor dependencies

$ docker-tools remote-build -g service/test:latest -i vendor
On branch develop
Your branch is up to date with 'origin/develop'.

nothing to commit, working tree clean
latest
fatal: ambiguous argument 'latest': unknown revision or path not in the working tree.
using the latest commit...
Executing: /bin/pwd [pwd]
/data/go/src/example.com/service/test
Executing: /usr/local/bin/git [git archive --format tar e8b54fa]
Executing: /usr/local/bin/gtar [gtar --transform s,^./,,g -rf /var/folders/ln/98ndp7416gx5v3mn7gv78ntm0000gn/T/1524649721-943271212 vendor]
Executing: /usr/bin/gzip [gzip -9f /var/folders/ln/98ndp7416gx5v3mn7gv78ntm0000gn/T/1524649721-943271212]
Executing: /bin/mv [mv /var/folders/ln/98ndp7416gx5v3mn7gv78ntm0000gn/T/1524649721-943271212.gz /var/folders/ln/98ndp7416gx5v3mn7gv78ntm0000gn/T/1524649721-943271212]
Executing: /usr/local/bin/gtar [gtar -tf /var/folders/ln/98ndp7416gx5v3mn7gv78ntm0000gn/T/1524649721-943271212]
.gitignore
.gitmodules
Dockerfile
...
Executing: mkdir -p /builds
Executing: /usr/bin/scp [scp -P 22 /var/folders/ln/98ndp7416gx5v3mn7gv78ntm0000gn/T/1524649721-943271212 root@127.0.0.1:/builds/1524649721-943271212.tar.gz]
1524649721-943271212                                                                          100%   13MB   7.4MB/s   00:01
Executing: mkdir -p /builds/1524649721-943271212
Executing: cd /builds/1524649721-943271212
Executing: tar -xf /builds/1524649721-943271212.tar.gz
Executing: docker build -f Dockerfile -t registry.example.com/service/test:e8b54fa .
Executing: docker push registry.example.com/service/test:e8b54fa
Executing: rm -f /builds/1524649721-943271212.tar.gz
Executing: rm -rf /builds/1524649721-943271212
Sending build context to Docker daemon   67.8MB
Step 1/11 : FROM registry.example.com/library/builder-go:1.9.2 as builder
 ---> 3b7a1b768e1d
Step 2/11 : COPY . /data/go/src/example.com/test
 ---> 1ddaf12d7713
...
 ---> Using cache
 ---> 2e762bf4e128
Successfully built 2e762bf4e128
Successfully tagged registry.example.com/service/test:e8b54fa
The push refers to repository [registry.example.com/service/test]
604d8539a0d8: Layer already exists
ff8f9e6fab20: Layer already exists
a4c0a34c75f3: Layer already exists
9dfa40a0da3b: Layer already exists
e8b54fa: digest: sha256:f51282c918fdb52e9ba5862f383d328ce559d03a1415611660623505f93e2893 size: 1157
```

### Import Docker Image

Subcommand `remote-import` of `docker-tools` pulls a given image on remote server and push the image to another registry.

All available options are:

``` Bash
$ docker-tools remote-import --help
...
Flags:
  -f, --from string   Registry to import image from (default "docker.io")
      --help          Show help
  -h, --host string   SSH host to run import procedures (default "127.0.0.1")
  -p, --port int      SSH port to connect (default 22)
  -t, --to string     Registry to export image to (default "registry.example.com")
  -u, --user string   SSH username (default "root")
...
```

Example:

``` Bash
$ docker-tools remote-import redis:4.0.2
Executing: docker pull docker.io/redis:4.0.2
Executing: docker tag docker.io/redis:4.0.2 registry.example.com/docker.io/redis:4.0.2
Executing: docker push registry.example.com/docker.io/redis:4.0.2
Executing: docker rmi docker.io/redis:4.0.2 registry.example.com/docker.io/redis:4.0.2
4.0.2: Pulling from library/redis
d13d02fa248d: Pull complete
039f8341839e: Pull complete
21b9cdda7eb9: Pull complete
c3eba3e5fbc2: Pull complete
7778a0753f87: Pull complete
b052cf77de81: Pull complete
Digest: sha256:cd277716dbff2c0211c8366687d275d2b53112fecbf9d6c86e9853edb0900956
Status: Downloaded newer image for redis:4.0.2
The push refers to repository [registry.example.com/docker.io/redis]
4aa04ab0fe76: Layer already exists
967b580842df: Layer already exists
22fc1222979f: Layer already exists
9503917b6420: Layer already exists
aa84bbcc6553: Layer already exists
29d71372a492: Layer already exists
4.0.2: digest: sha256:3c07847e5aa6911cf5d9441642769d3b6cd0bf6b8576773ae3a0742056b9dd47 size: 1571
Untagged: redis:4.0.2
Untagged: redis@sha256:cd277716dbff2c0211c8366687d275d2b53112fecbf9d6c86e9853edb0900956
Untagged: registry.example.com/docker.io/redis:4.0.2
Untagged: registry.example.com/docker.io/redis@sha256:3c07847e5aa6911cf5d9441642769d3b6cd0bf6b8576773ae3a0742056b9dd47
Deleted: sha256:8f2e175b3bd129fd9416df32a0e51f36632e3ab82c5608b4030590ad79f0be12
Deleted: sha256:dc220825bf188145846d269e04e122f8d53194a8b18652df23410ba114dde020
Deleted: sha256:f76a80a6c86476894da4c51e6415cc0201dbcd75b135e74e069d09bc51bcd094
Deleted: sha256:dbe291a244f66ac9ab2d76f3106e21f23479ab966e8437a7dba9ac13b0a9a793
Deleted: sha256:f9a111ff6d25e72d0448c4490ee6f3296ce44653933951fe96e50dcd809f35dc
Deleted: sha256:72ed9453809fb6a19a7a5101af4efbf1328c87b1d24ce02cc73cfd034f125166
Deleted: sha256:29d71372a4920ec230739a9e2317e7e9b18644edb10f78cde85df85e6ab85fc2
```

## License

The Apache License (Version 2.0, January 2004).
