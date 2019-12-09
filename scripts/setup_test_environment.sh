#!/bin/bash -ex

HELM_V2_VERSION="2.16.1"
HELM_V3_VERSION="3.0.1"
CHARTMUSEUM_VERSION="0.8.2"

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR/../

export PATH="$PWD/testbin:$PATH"
export TEST_HELM_HOME="$PWD/.helm"

[ "$(uname)" == "Darwin" ] && PLATFORM="darwin" || PLATFORM="linux"

main() {
    install_helm_v2
    install_helm_v3
    install_chartmuseum
    package_test_charts
}

install_helm_v2() {
    if [ ! -f "testbin/helm2" ]; then
        mkdir -p testbin/
        TARBALL="helm-v${HELM_V2_VERSION}-${PLATFORM}-amd64.tar.gz"
        wget "https://get.helm.sh/${TARBALL}"
        tar -C testbin/ -xzf $TARBALL
        rm -f $TARBALL
        pushd testbin/
        UNCOMPRESSED_DIR="$(find . -mindepth 1 -maxdepth 1 -type d)"
        mv $UNCOMPRESSED_DIR/helm .
        rm -rf $UNCOMPRESSED_DIR
        chmod +x ./helm
        mv ./helm ./helm2
        popd
        HELM_HOME=${TEST_HELM_HOME} helm2 init --client-only
    fi
}

install_helm_v3() {
    if [ ! -f "testbin/helm3" ]; then
        mkdir -p testbin/
        TARBALL="helm-v${HELM_V3_VERSION}-${PLATFORM}-amd64.tar.gz"
        wget "https://get.helm.sh/${TARBALL}"
        tar -C testbin/ -xzf $TARBALL
        rm -f $TARBALL
        pushd testbin/
        UNCOMPRESSED_DIR="$(find . -mindepth 1 -maxdepth 1 -type d)"
        mv $UNCOMPRESSED_DIR/helm .
        rm -rf $UNCOMPRESSED_DIR
        chmod +x ./helm
        mv ./helm ./helm3
        popd
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
    pushd testdata/charts/helm2/
    for d in $(find . -maxdepth 1 -mindepth 1 -type d); do
        pushd $d
        HELM_HOME=${TEST_HELM_HOME} helm2 package --sign --key helm-test --keyring ../../../pgp/helm-test-key.secret .
        popd
    done
    popd

    pushd testdata/charts/helm3/
    for d in $(find . -maxdepth 1 -mindepth 1 -type d); do
        pushd $d
        helm3 package --sign --key helm-test --keyring ../../../pgp/helm-test-key.secret .
        popd
    done
    popd
}

main
