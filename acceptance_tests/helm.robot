*** Settings ***
Documentation     Tests to verify that helm-push can be installed
...               and run as a Helm plugin etc.
Library           lib/Helm.py
Suite Setup       Suite Setup
Suite Teardown    Suite Teardown

*** Test Cases ***
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

*** Keywords ***
Suite Setup
    remove helm plugin

Suite Teardown
    remove helm plugin
