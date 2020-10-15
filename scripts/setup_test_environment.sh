#!/bin/bash -ex

HELM_V2_VERSION="v2.16.12"
HELM_V3_VERSION="v3.3.4"
CHARTMUSEUM_VERSION="v0.12.0"

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR/../

export PATH="$PWD/testbin:$PATH"
export TEST_V2_HELM_HOME="$PWD/.helm2"
export TEST_V3_XDG_CACHE_HOME="$PWD/.helm3/xdg/cache"
export TEST_V3_XDG_CONFIG_HOME="$PWD/.helm3/xdg/config"
export TEST_V3_XDG_DATA_HOME="$PWD/.helm3/xdg/data"

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
        TARBALL="helm-${HELM_V2_VERSION}-${PLATFORM}-amd64.tar.gz"
        curl -LO "https://get.helm.sh/${TARBALL}"
        tar -C testbin/ -xzf $TARBALL
        rm -f $TARBALL
        pushd testbin/
        UNCOMPRESSED_DIR="$(find . -mindepth 1 -maxdepth 1 -type d)"
        mv $UNCOMPRESSED_DIR/helm .
        rm -rf $UNCOMPRESSED_DIR
        chmod +x ./helm
        mv ./helm ./helm2
        popd
        HELM_HOME=${TEST_V2_HELM_HOME} helm2 init --client-only
    fi
}

install_helm_v3() {
    if [ ! -f "testbin/helm3" ]; then
        mkdir -p testbin/
        TARBALL="helm-${HELM_V3_VERSION}-${PLATFORM}-amd64.tar.gz"
        curl -LO "https://get.helm.sh/${TARBALL}"
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
        curl -LO "https://s3.amazonaws.com/chartmuseum/release/${CHARTMUSEUM_VERSION}/bin/${PLATFORM}/amd64/chartmuseum"
        chmod +x ./chartmuseum
        popd
    fi
}

package_test_charts() {
    pushd testdata/charts/helm2/
    for d in $(find . -maxdepth 1 -mindepth 1 -type d); do
        pushd $d
        HELM_HOME=${TEST_V2_HELM_HOME} helm2 package --sign --key helm-test --keyring ../../../pgp/helm-test-key.secret .
        popd
    done
    popd

    pushd testdata/charts/helm3/
    for d in $(find . -maxdepth 1 -mindepth 1 -type d); do
        pushd $d
        XDG_CACHE_HOME=${TEST_V3_XDG_CACHE_HOME} XDG_CONFIG_HOME=${TEST_V3_XDG_CONFIG_HOME} XDG_DATA_HOME=${TEST_V3_XDG_DATA_HOME} helm3 package \
            --sign --key helm-test --keyring ../../../pgp/helm-test-key.secret .
        popd
    done
    popd
}

main
