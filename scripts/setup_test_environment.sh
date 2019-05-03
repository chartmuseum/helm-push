#!/bin/bash -ex

HELM_VERSION="2.13.1"
CHARTMUSEUM_VERSION="0.8.2"

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR/../

export PATH="$PWD/testbin:$PATH"
export HELM_HOME="$PWD/.helm"

[ "$(uname)" == "Darwin" ] && PLATFORM="darwin" || PLATFORM="linux"

main() {
    install_helm
    install_chartmuseum
    package_test_charts
}

install_helm() {
    if [ ! -f "testbin/helm" ]; then
        mkdir -p testbin/
        TARBALL="helm-v${HELM_VERSION}-${PLATFORM}-amd64.tar.gz"
        wget "https://storage.googleapis.com/kubernetes-helm/${TARBALL}"
        tar -C testbin/ -xzf $TARBALL
        rm -f $TARBALL
        pushd testbin/
        UNCOMPRESSED_DIR="$(find . -mindepth 1 -maxdepth 1 -type d)"
        mv $UNCOMPRESSED_DIR/helm .
        rm -rf $UNCOMPRESSED_DIR
        chmod +x ./helm
        popd
        helm init --client-only
    fi
}

install_chartmuseum() {
    if [ ! -f "testbin/chartmuseum" ]; then
        mkdir -p testbin/
        pushd testbin/
        wget "https://s3.amazonaws.com/chartmuseum/release/v${CHARTMUSEUM_VERSION}/bin/${PLATFORM}/amd64/chartmuseum"
        chmod +x ./chartmuseum
        popd
    fi
}

package_test_charts() {
    pushd testdata/charts/
    for d in $(find . -maxdepth 1 -mindepth 1 -type d); do
        pushd $d
        helm package --sign --key helm-test --keyring ../../pgp/helm-test-key.secret .
        popd
    done
    popd
}

main
