---
name: mapping-os-packages-to-chainguard
description: Instructions for mapping the name of OS packages from ecosystems like Alpine, Debian and Red Hat to Chainguard. Use this when converting Dockerfiles to Chainguard.
---

These are the steps you should follow when trying to map an OS package name
from a distribution like Debian or Alpine to its equivalent in the Chainguard
ecosystem. 

At a high level this involves these steps:

1. Download package mappings
2. Map package name
3. Optionally, search for packages

## Download Package Mappings

The Dockerfile Converter tool provides a list of package mappings that can be
useful for this task. You can download them like this:

```
curl -sSf https://raw.githubusercontent.com/chainguard-dev/dfc/refs/heads/main/pkg/dfc/builtin-mappings.yaml > builtin-mappings.yaml 
```

The structure of this file is like:

```
# You can ignore this for the purposes of image mapping
images: {}

packages:
    alpine:
        alpine-package-name:
            - chainguard-equivalent-name
    debian:
        debian-package-name:
            - chainguard-equivalent-name
    fedora:
        fedora-package-name:
            - chainguard-equivalent-name
```

## Map Package Names

You can use `yq` to search the mappings for the the equivalent Chainguard
packages for a given package in another distribution: 

```
yq '.packages.debian["awscli"]' builtin-mappings.yaml
```

For Debian based distributions (like Debian and Ubunutu) use `.packages.debian`.

For Alpine, use `.packages.alpine`.

For Red Hat and Fedora, use `.packages.fedora`.

If you can't find a match, then assume that the package name is the same in the
Chainguard ecosystem, which is often the case.

## Search For Packages

If you find that a package name doesn't exist in Chainguard's repositories, you
can drop into a `-dev` image and use `apk search` to find the equivalent based
on naming:

```
docker run -it --rm --entrypoint bash -u root cgr.dev/ORGANIZATION/python:3.12-dev

# apk update
# apk search -q <package name or substring>
# apk search -q cmd:<specific command name>
```

Another thing you may consider trying is seeing which files are provided by the
upstream:

```
docker run -it --rm debian:bookworm-slim

# apt update
# apt install mariadb-client
# dpkg -L mariadb-client
```

And then you can try and find packages that provide those files:

```
docker run -it --rm --entrypoint bash -u root cgr.dev/ORGANIZATION/python:3.12-dev

# apk update
# apk search -q cmd:mysqlshow
```
