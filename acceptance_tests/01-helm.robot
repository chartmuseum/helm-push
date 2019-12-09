*** Settings ***
Documentation     Tests to verify that helm-push can be installed
...               and run as a Helm plugin etc.
Library           lib/Helm.py
Suite Setup       Suite Setup
Suite Teardown    Suite Teardown

*** Test Cases ***
Plugin installs on Helm 2
    Test plugin installation   2

Plugin installs on Helm 3
    Test plugin installation   3

*** Keywords ***
Test plugin installation
    [Arguments]    ${version}
    set helm version    ${version}
    helm-push can be installed as a Helm plugin
    helm-push is listed as a Helm plugin after install
    helm-push can be run as a Helm plugin
    helm-push can be removed
    helm-push is not listed as a Helm plugin after removal

helm-push can be installed as a Helm plugin
    install helm plugin
    return code should be   0

helm-push is listed as a Helm plugin after install
    check helm plugin
    return code should be   0

helm-push can be run as a Helm plugin
    run helm plugin
    return code should be   0

helm-push can be removed
    remove helm plugin
    return code should be   0

helm-push is not listed as a Helm plugin after removal
    check helm plugin
    return code should not be   0

Suite Setup
    set helm version    3
    remove helm plugin
    set helm version    2
    remove helm plugin

Suite Teardown
    set helm version    3
    remove helm plugin
    set helm version    2
    remove helm plugin
