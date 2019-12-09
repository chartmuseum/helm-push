package helm

import (
	"os"
	"os/exec"
)

type (
	HelmMajorVersion int
)

const (
	HelmMajorVersion2 = 2
	HelmMajorVersion3 = 3
)

func GetHelmMajorVersion() HelmMajorVersion {
	helmBin, helmBinVarSet := os.LookupEnv("HELM_BIN")
	if !helmBinVarSet {
		helmBin = "helm"
	}
	helmVersion2CheckCmd := exec.Command(helmBin, "version", "-c", "--tls")
	err := helmVersion2CheckCmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return HelmMajorVersion3
	}
	return HelmMajorVersion2
}