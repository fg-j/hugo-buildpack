# Hugo Cloud Native Buildpack

The Hugo CNB provides [Hugo](https://gohugo.io) and uses it to generate static
assets. These static assets can be served with a server like HTTPD. The
buildpack does *not* install any server for the assets.

The buildpack detects true if the app contains a `content` directory that
contains at least one Markdown or HTML file.

The buildpack places generated static files in `<app root>/output`.

## Integration

The Hugo CNB provides the `hugo` dependency (it puts the Hugo CLI on the
`$PATH`). Downstream buildpacks can require the `hugo` dependency by generating
a [Build Plan
TOML](https://github.com/buildpacks/spec/blob/master/buildpack.md#build-plan-toml)
file that looks like the following:

```toml
[[requires]]

  # The name of the  dependency is "hugo".
  name = "hugo"

  # The Hugo buildpack supports some non-required metadata options.
  [requires.metadata]

    # Setting the build flag to true will ensure that the Hugo
    # depdendency is available on the $PATH for subsequent buildpacks during
    # their build phase. If you are writing a buildpack that needs to run Hugo
    # during its build process, this flag should be set to true.
    build = true

    # Setting the launch flag to true will ensure that the Hugo
    # dependency is available on the $PATH for the running application. If you are
    # writing an application that needs to run Hugo at runtime, this flag should
    # be set to true.
    launch = true
```

## Usage

To package this buildpack for consumption:

```
$ ./scripts/package.sh --version <x.y.z>
```

This builds the buildpack's Go source using `GOOS=linux` by default. You can
supply another value as the first argument to `package.sh`.

Once you have packaged the buildpack, you can build a site with it using the
[pack CLI](https://github.com/buildpacks/pack):
```bash
cd my-hugo-site-directory
pack build hugo-site /
--buildpack paketobuildpacks/httpd:latest
--buildpack /path/to/hugo-buildpack/build/buildpack.tgz
--buildpack gcr.io/paketo-community/build-plan:latest
--builder paketobuildpacks/builder:full
```

Note: The above example assumes that you'll use
[HTTPD](https://httpd.apache.org/) to serve the static assets. The [Paketo
HTTPD CNB](https://github.com/paketo-buildpacks/httpd) can install the server
in your app container. See the sample app in
[`integration/testdata/default_site`](integration/testdata/default_site) for a
complete working example app.
