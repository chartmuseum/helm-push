# helm push plugin
<img align="right" src="https://github.com/kubernetes-helm/chartmuseum/raw/master/logo.png">

[![Codefresh build status]( https://g.codefresh.io/api/badges/build?repoOwner=chartmuseum&repoName=helm-push&branch=master&pipelineName=helm-push&accountName=codefresh-inc&type=cf-2)]( https://g.codefresh.io/repositories/chartmuseum/helm-push/builds?filter=trigger:build;branch:master;service:5ad4eed637adc30001207fab~helm-push)

Helm plugin to push chart package to [ChartMuseum](https://github.com/kubernetes-helm/chartmuseum)

## Install
Based on the version in `plugin.yaml`, release binary will be downloaded from GitHub:

```
$ helm plugin install https://github.com/chartmuseum/helm-push
Downloading and installing helm-push v0.1.0 ...
https://github.com/chartmuseum/helm-push/releases/download/v0.1.0/helm-push_0.1.0_darwin_amd64.tar.gz
Installed plugin: push
```

## Usage
You must start by adding a ChartMuseum-backed repo via Helm CLI (if not already added)
```
$ helm repo add chartmuseum http://localhost:8080
```
For all available plugin options, please run
```
$ helm push --help
```

### Pushing a directory
Point to a directory containing a valid `Chart.yaml` and the chart will be packaged and uploaded:
```
$ cat mychart/Chart.yaml
name: mychart
version: 0.3.2
```
```
$ helm push mychart/ chartmuseum
Pushing mychart-0.3.2.tgz to chartmuseum...
Done.
```

### Pushing with a custom version
The `--version` flag can be provided, which will push the package with a custom version.

Here is an example using the last git commit id as the version:
```
$ helm push mychart/ --version="$(git log -1 --pretty=format:%h)" chartmuseum
Pushing mychart-5abbbf28.tgz to chartmuseum...
Done.
```
If you want to enable something like `--version="latest"`, which you intend to push regularly, you will need to run your ChartMuseum server with `ALLOW_OVERWRITE=true`.

### Push .tgz package
This workflow does not require the use of `helm package`, but pushing .tgzs is still suppported:
```
$ helm push mychart-0.3.2.tgz chartmuseum
Pushing mychart-0.3.2.tgz to chartmuseum...
Done.
```

## Authentication
### Basic Auth
If you have added your repo with the `--username`/`--password` flags (Helm 2.9+), or have added your repo with the basic auth username/password in the URL (e.g. `https://myuser:mypass@my.chart.repo.com`), no further setup is required.

The plugin will use the auth info located in `~/.helm/repository/repositories.yaml` in order to authenticate.

If you are running ChartMuseum with `AUTH_ANONYMOUS_GET=true`, and have added your repo without authentication, the plugin recognizes the following environment variables for basic auth on push operations:
```
$ export HELM_REPO_USERNAME="myuser"
$ export HELM_REPO_PASSWORD="mypass"
```

With this setup, you can enable people to use your repo for installing charts etc. without allowing them to upload to it.

### Token
Although ChartMuseum server does not define or accept a token format (yet), if you are running it behind a proxy that accepts access tokens, you can provide the following env var:
```
$ export HELM_REPO_ACCESS_TOKEN="<token>"
```

This will result in all basic auth options above being ignored, and the plugin will send the token in the header:
```
Authorization: Bearer <token>
```
