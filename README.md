# helm cm-push plugin
<img align="right" src="https://github.com/helm/chartmuseum/raw/main/logo.png">

[![GitHub Actions status](https://github.com/chartmuseum/helm-push/workflows/build/badge.svg)](https://github.com/chartmuseum/helm-push/actions?query=workflow%3Abuild)

Helm plugin to push chart package to [ChartMuseum](https://github.com/helm/chartmuseum)

## Install
Based on the version in `plugin.yaml`, release binary will be downloaded from GitHub:

```
$ helm plugin install https://github.com/chartmuseum/helm-push.git
Downloading and installing helm-push v0.9.0 ...
https://github.com/chartmuseum/helm-push/releases/download/v0.9.0/helm-push_0.9.0_darwin_amd64.tar.gz
Installed plugin: cm-push
```

## Usage
Start by adding a ChartMuseum-backed repo via Helm CLI (if not already added)
```
$ helm repo add chartmuseum http://localhost:8080
```
For all available plugin options, please run
```
$ helm cm-push --help
```

### Pushing a directory
Point to a directory containing a valid `Chart.yaml` and the chart will be packaged and uploaded:
```
$ cat mychart/Chart.yaml
name: mychart
version: 0.3.2
```
```
$ helm cm-push mychart/ chartmuseum
Pushing mychart-0.3.2.tgz to chartmuseum...
Done.
```

### Pushing with a custom version
The `--version` flag can be provided, which will push the package with a custom version.

Here is an example using the last git commit id as the version:
```
$ helm cm-push mychart/ --version="$(git log -1 --pretty=format:%h)" chartmuseum
Pushing mychart-5abbbf28.tgz to chartmuseum...
Done.
```
If you want to enable something like `--version="9.9.9-dev1"`, which you intend to push regularly, you will need to run your ChartMuseum server with `ALLOW_OVERWRITE=true`.

### Push .tgz package
This workflow does not require the use of `helm package`, but pushing .tgzs is still suppported:
```
$ helm cm-push mychart-0.3.2.tgz chartmuseum
Pushing mychart-0.3.2.tgz to chartmuseum...
Done.
```

### Force push
If your ChartMuseum install is configured with `ALLOW_OVERWRITE=true`, chart versions will be automatically overwritten upon re-upload.

Otherwise, unless your install is configured with `DISABLE_FORCE_OVERWRITE=true` (ChartMuseum > v0.7.1), you can use the `--force`/`-f` option to to force an upload:
```
$ helm cm-push --force mychart-0.3.2.tgz chartmuseum
Pushing mychart-0.3.2.tgz to chartmuseum...
Done.
```

### Pushing directly to URL
If the second argument provided resembles a URL, you are not required to add the repo prior to push:
```
$ helm cm-push mychart-0.3.2.tgz http://localhost:8080
Pushing mychart-0.3.2.tgz to http://localhost:8080...
Done.
```

## Context Path

If you are running ChartMuseum behind a proxy that adds a route prefix, for example:
```
https://my.chart.repo.com/helm/v1/index.yaml -> http://chartmuseum-svc/index.yaml
```

You can use the `--context-path=` option or `HELM_REPO_CONTEXT_PATH` env var in order for the plugin to construct the upload URL correctly:
```
helm repo add chartmuseum https://my.chart.repo.com/helm/v1
helm cm-push --context-path=/helm/v1 mychart-0.3.2.tgz chartmuseum
```

Alternatively, you can add `serverInfo.contextPath` to your index.yaml:
```
apiVersion: v1
entries:{}
generated: "2018-08-09T11:08:21-05:00"
serverInfo:
  contextPath: /helm/v1
```

In ChartMuseum server (>0.7.1) this will automatically be added to index.yaml if the `--context-path` option is provided.

## Authentication
### Basic Auth
If you have added your repo with the `--username`/`--password` flags (Helm 2.9+), or have added your repo with the basic auth username/password in the URL (e.g. `https://myuser:mypass@my.chart.repo.com`), no further setup is required.

The plugin will use the auth info located in `~/.helm/repository/repositories.yaml` (for Helm 2) or `~/.config/helm/repositories.yaml` (for Helm 3) in order to authenticate.

If you are running ChartMuseum with `AUTH_ANONYMOUS_GET=true`, and have added your repo without authentication, the plugin recognizes the following environment variables for basic auth on push operations:
```
$ export HELM_REPO_USERNAME="myuser"
$ export HELM_REPO_PASSWORD="mypass"
```

With this setup, you can enable people to use your repo for installing charts etc. without allowing them to upload to it.

### Token

*ChartMuseum token-auth is currently in progress. Pleasee see [auth-server-example](https://github.com/chartmuseum/auth-server-example) for more info.*

Although ChartMuseum server does not define or accept a token format (yet), if you are running it behind a proxy that accepts access tokens, you can provide the following env var:
```
$ export HELM_REPO_ACCESS_TOKEN="<token>"
```

This will result in all basic auth options above being ignored, and the plugin will send the token in the header:
```
Authorization: Bearer <token>
```

If you require a custom header to be used for passing the token, you can the following env var:
```
$ export HELM_REPO_AUTH_HEADER="<myheader>"
```

This will then be used in place of `Authorization: Bearer`:
```
<myheader>: <token>
```

#### Token config file (~/.cfconfig)
For users of [Managed Helm Repositories](https://codefresh.io/codefresh-news/introducing-managed-helm-repositories/) (Codefresh), the plugin is able to auto-detect your API key from `~/.cfconfig`. This file is managed by [Codefresh CLI](https://codefresh-io.github.io/cli/).

If detected, this API key will be used for token-based auth, overriding basic auth options described above.

The format of this file is the following:

```
contexts:
  default:
    name: default
    token: <token>
current-context: default
```

### TLS Client Cert Auth

ChartMuseum server does not yet have options to setup TLS client cert authentication (please see [chartmuseum#79](https://github.com/helm/chartmuseum/issues/79)).

If you are running ChartMuseum behind a frontend that does, the following options are available:

```
--ca-file string    Verify certificates of HTTPS-enabled servers using this CA bundle [$HELM_REPO_CA_FILE]
--cert-file string  Identify HTTPS client using this SSL certificate file [$HELM_REPO_CERT_FILE]
--key-file string   Identify HTTPS client using this SSL key file [$HELM_REPO_KEY_FILE]
--insecure          Connect to server with an insecure way by skipping certificate verification [$HELM_REPO_INSECURE]
```

## Custom Downloader
This plugin also defines the `cm://` protocol that you may specify when adding a repo:
```
$ helm repo add chartmuseum cm://my.chart.repo.com
```

The only real difference with this vs. simply using http/https, is that the environment variables above are recognized by the plugin and used to set the `Authorization` header appropriately. As in, if you do not add your repo in this way, you are unable to use token-based auth for GET requests (downloading index.yaml, chart .tgzs, etc).

By default, `cm://` translates to `https://`. If you must use `http://`, you can set the following env var:
```
$ export HELM_REPO_USE_HTTP="true"
```
