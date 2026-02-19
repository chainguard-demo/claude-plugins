---
name: migrating-dockerfiles-to-chainguard
description: Instructions for converting existing Dockerfiles to use Chainguard images.
---

To convert a Dockerfile to use Chainguard images you should follow these steps:

1. Clarify instructions with the user
2. Convert the Dockerfile
   1. Modify `FROM` statements to use Chainguard images.
   2. Migrate any uses of other OS package managers to `apk`
   3. Adjust users and permissions where required.
3. Test the conversion
   1. Build the image
   2. Run the image
   3. Troubleshoot any failures

## Clarify Instructions

If you're asked to convert a Dockerfile to Chainguard, before doing anything
else, you should clarify the following with the user: 

1. What is the Chainguard organization name? Each Chainguard customer has their
   own organization and it will be referred to generically as `ORGANIZATION` in
   these instructions.
2. Is FIPS required? Chainguard offers FIPS compliant images that have the
   suffix `-fips`. For instance: `python-fips`.
3. Should images be pulled directly from `cgr.dev`? Or, if images are
   hosted/proxied from another repository, what is that repository?
4. Should you test the conversion?
   1. Should you try and build the image to validate the conversion worked?
   2. Should you try and run the image to validate the conversion worked?

Store these answers under `~/.claude/chainguard-preferences.md` and try and reuse
this file across subsequent asks to migrate Dockerfiles.

## Backup the Dockerfile

You should backup the old Dockerfile by prepending `old.` to it. For instance,
`Dockerfile` -> `old.Dockerfile`.

Then you can modify the original file directly.

## Converting FROM statements

You should replace images defined in `FROM` lines with the equivalent
Chainguard image.

Always use the `mapping-container-images-to-chainguard` skill if it is
available.

You should try to match to actively maintained tags when converting Dockerfiles.

You should prefer `-dev` tags as they have more utilities in them and they are
less likely to cause issues. One known exception is `chainguard-base` which has
no `-dev` variant. You should always use the `latest` tag for `chainguard-base`.

### Digests

If the existing reference uses a digest reference (i.e
`python:3.12-slim-bookworm@sha256:a866731a6b71c4a194a845d86e06568725e430ed21821d0c52e4efb385cf6c6f`)
then you should also include a digest. Naturally, the digest will be different
for the Chainguard image. You can figure out the digest in a few ways, depending
on the available tooling:

```
# With crane (prefer this if it is available)
crane digest cgr.dev/ORGANIZATION/python:3.12-dev

# Or, with docker
docker pull cgr.dev/ORGANIZATION/python:3.12-dev
docker inspect --format='{{index .RepoDigests 0}}' cgr.dev/ORGANIZATION/python:3.12-dev
```

### Chainguard Base

Chainguard have their own Linux distribution, so they don't have equivalent
images for `alpine`, `debian`, `ubuntu` etc. 

Where the Dockerfile is using one of these generic bases, you should swap it
to use `chainguard-base` (or `chainguard-base-fips` for FIPS). You should
always use the `latest` tag with `chainguard-base`. There is no `latest-dev`
tag for `chainguard-base`, so don't use that.

## Multi Stage Builds

### Dev Variants

Where there are multiple `FROM` lines in a Dockerfile, you should try to use a
non-dev variant (i.e `cgr.dev/ORGANIZATION/python:3.12`) for the final stage
unless there are any `RUN` lines that follow it (in which case you need a shell,
and therefore a `-dev` image).

### COPY Permissions

A common issue with multi stage builds is that files are created in the first
stage by the `root` user, which are then not accessible to the `65532` user in
the final stage. You can mitigate this by ensuring files copied into the final
stage belong to the `65532` user:

```
COPY --from=build --chown=65532:65532 /app .
```

## Adding Packages with `apk`

Unlike other Linux distributions that may use `apt`, `yum` or `dnf` to install
packages, Chainguard images use `apk` to manage packages.

This means that you should convert instances of `apt`, `yum` and `dnf` to `apk`
in `RUN` lines.

For instance, for `apt`, lines like this:

```
RUN apt-get update \
    && apt-get install -y software-properties-common=0.99.22.9 \
    && add-apt-repository ppa:libreoffice/libreoffice-still \
    && apt-get install -y libreoffice \
    && apt-get clean
```

Should become:

```
RUN apk --no-cache add libreoffice
```

And lines like this:

```
RUN yum update -y && \
    yum -y install httpd php php-cli php-common && \
    yum clean all && \
    rm -rf /var/cache/yum/*

RUN dnf update -y && \
    dnf -y install httpd php php-cli php-common && \
    dnf clean all

RUN microdnf update -y && \
    microdnf -y install httpd php php-cli php-common && \
    microdnf clean all
```

Should become:

```
RUN apk add --no-cache httpd php php-cli php-common
```

Don't just remove `apt` lines. Replace them with `apk` equivalents.

### Root Permissions

Chainguard images typically run as a non root user with the UID `65532`, rather
than `root`. You will need to switch to `root` before running `apk add` and
then revert back to `65532` after.

```
USER root
RUN apk add --no-cache httpd php php-cli php-common
USER 65532
```

## Users

As mentioned above, Chainguard images typically run as a non root user with the
UID `65532`, rather than `root`. Commonly, this user is called `nonroot`, but
not always. For instance, in the `node` image the user is called `node`.

To be safe, always prefer the UID when issuing `USER` statements or changing
the permissions of files.

You should switch to the `root` user with `USER root` before running privileged
commands (like `npm`) and then use `USER 65532` after to return to the non root
user.

You can figure out which user the image is configured to use by querying the
image configuration. How you do that depends on the tools you have available:

```
# With crane.
crane config cgr.dev/chainguard/python | jq -r '.config.User'

# With docker.
docker inspect --format='{{ .Config.User }}' cgr.dev/chainguard/python
```

You can also find this information in the Specifications tab of the Chainguard
Image Directory:

```
https://images.chainguard.dev/directory/image/python/specifications
```

## Entrypoints

Unlike upstream images, the entrypoint for Chainguard base images is typically
the runtime interpreter for the given language (i.e `java`, `python`, `php`)
rather than a shell like `sh` or `bash`.

This means that a `CMD` like `CMD ["python", "app.py"]` will cause an error in a
Chainguard images because it translates to running `python python app.py`.

We should prefer to use explicit entrypoints like:

```
ENTRYPOINT ["python", "app.py"]
```

## PHP

It's common in Dockerfiles that build PHP applications to install composer in
this way:

```
COPY --from=composer:latest /usr/bin/composer /usr/local/bin/composer
```

This is unnecessary when using the Chainguard `php` image because `composer` is
already included in the `-dev` tags. Therefore, you can remove any lines that
install composer.

## Troubleshooting

### Missing Shell

If you get errors complaining about a missing `sh` or `bash`, then ensure you
are using a `-dev` image, which provides a shell.

### Permission Denied

If you get errors complaining about permissions, then it's probably to
do with the non root user. Switch to `USER root` temporarily and then switch
back to `USER 65532` once you've completed the privileged operation.

### Check the Overview page

The Overview page for the Chainguard image typically includes information that
is helpful when performing a migration.

```
https://images.chainguard.dev/directory/image/<image-name>/overview
```

If you run into issues, it is worth checking this page for guidance.
