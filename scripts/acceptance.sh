#!/bin/bash -ex

PY_REQUIRES="robotframework==4.1.1"

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR/../

if [ "$(uname)" == "Darwin" ]; then
    PLATFORM="darwin"
else
    PLATFORM="linux"
fi

export PATH="$PWD/testbin:$PWD/bin/$PLATFORM/amd64:$PATH"
export HELM_PUSH_PLUGIN_NO_INSTALL_HOOK=1

export TEST_V2_HELM_HOME="$PWD/.helm2"
HELM_HOME=${TEST_V2_HELM_HOME} helm2 init --client-only

export TEST_V3_XDG_CACHE_HOME="$PWD/.helm3/xdg/cache"
export TEST_V3_XDG_CONFIG_HOME="$PWD/.helm3/xdg/config"
export TEST_V3_XDG_DATA_HOME="$PWD/.helm3/xdg/data"

if [ ! -d .venv/ ]; then
    virtualenv -p $(which python3) .venv/
    .venv/bin/python .venv/bin/pip3 install $PY_REQUIRES
fi

mkdir -p .robot/
.venv/bin/robot --outputdir=.robot/ acceptance_tests/
