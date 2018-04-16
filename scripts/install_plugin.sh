#!/bin/bash -e

# Copied w/ love from the excellent hypnoglow/helm-s3

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR/../

if [ -n "${HELM_PUSH_PLUGIN_NO_INSTALL_HOOK}" ]; then
    echo "Development mode: not downloading versioned release."
    exit 0
fi

version="$(cat plugin.yaml | grep "version" | cut -d '"' -f 2)"
echo "Downloading and installing helm-push v${version} ..."

url=""
if [ "$(uname)" == "Darwin" ]; then
    url="https://github.com/chartmuseum/helm-push/releases/download/v${version}/helm-push_${version}_darwin_amd64.tar.gz"
elif [ "$(uname)" == "Linux" ] ; then
    url="https://github.com/chartmuseum/helm-push/releases/download/v${version}/helm-push_${version}_linux_amd64.tar.gz"
else
    url="https://github.com/chartmuseum/helm-push/releases/download/v${version}/helm-push_${version}_windows_amd64.tar.gz"
fi

echo $url

mkdir -p "bin"
mkdir -p "releases/v${version}"

# Download with curl if possible.
if [ -x "$(which curl 2>/dev/null)" ]; then
    curl -sSL "${url}" -o "releases/v${version}.tar.gz"
else
    wget -q "${url}" -O "releases/v${version}.tar.gz"
fi
tar xzf "releases/v${version}.tar.gz" -C "releases/v${version}"
mv "releases/v${version}/bin/helmpush" "bin/helmpush" || \
    mv "releases/v${version}/bin/helmpush.exe" "bin/helmpush"

