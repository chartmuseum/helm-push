*** Settings ***
Documentation     Tests to verify that helm-push can be used to
...               successfully push a package to chartmuseum server
Library           lib/ChartMuseum.py
Library           lib/Helm.py
Library           lib/HelmPush.py
Suite Setup       Suite Setup
Suite Teardown    Suite Teardown

*** Test Cases ***
Chart directory can be pushed to ChartMuseum
    push chart directory
    HelmPush.return code should be   0
    package exists in chartmuseum storage
    ChartMuseum.return code should be   0
    clear chartmuseum storage

Chart directory can be pushed to ChartMuseum with custom version
    push chart directory    latest
    HelmPush.return code should be   0
    package exists in chartmuseum storage   latest
    ChartMuseum.return code should be   0
    clear chartmuseum storage

Chart package can be pushed to ChartMuseum
    push chart package
    HelmPush.return code should be   0
    package exists in chartmuseum storage
    ChartMuseum.return code should be   0
    clear chartmuseum storage

Chart package can be pushed to ChartMuseum with custom version
    push chart package  latest
    HelmPush.return code should be   0
    package exists in chartmuseum storage   latest
    ChartMuseum.return code should be   0
    clear chartmuseum storage

*** Keywords ***
Suite Setup
    remove chartmuseum logs
    start chartmuseum
    Sleep  2
    add chart repo

Suite Teardown
    remove chart repo
    stop chartmuseum
    print chartmuseum logs
