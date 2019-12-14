module github.com/chartmuseum/helm-push

go 1.13

require (
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/spf13/cobra v0.0.5
	helm.sh/helm/v3 v3.0.1
	k8s.io/helm v2.16.1+incompatible
)

replace github.com/docker/docker => github.com/docker/docker v0.0.0-20190731150326-928381b2215c
