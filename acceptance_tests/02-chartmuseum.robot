*** Settings ***
Documentation     Tests to verify that helm-push can be used to
...               successfully push a package to chartmuseum server
Library           lib/ChartMuseum.py
Library           lib/Helm.py
Library           lib/HelmPush.py
Suite Setup       Suite Setup
Suite Teardown    Suite Teardown

*** Test Cases ***
Plugin works with ChartMuseum on Helm 2
    Test ChartMuseum integration   2

Plugin works with ChartMuseum on Helm 3
    Test ChartMuseum integration   3

*** Keywords ***
Test ChartMuseum integration
    [Arguments]    ${version}
    set helm version    ${version}
    install helm plugin
    helm major version detected by plugin is  ${version}
    clear chartmuseum storage
    Chart directory can be pushed to ChartMuseum
    Chart directory can be pushed to ChartMuseum with custom version
    Chart package can be pushed to ChartMuseum
    Chart package can be pushed to ChartMuseum with custom version

Chart directory can be pushed to ChartMuseum
    # Repo name
    push chart directory
    HelmPush.return code should be   0
    package exists in chartmuseum storage
    package contains expected files
    HelmPush.return code should be   0
    ChartMuseum.return code should be   0
    clear chartmuseum storage

    # Repo URL
    push chart directory to url
    HelmPush.return code should be   0
    package exists in chartmuseum storage
    package contains expected files
    HelmPush.return code should be   0
    ChartMuseum.return code should be   0
    clear chartmuseum storage

Chart directory can be pushed to ChartMuseum with custom version
    # Repo name
    push chart directory    latest
    HelmPush.return code should be   0
    package exists in chartmuseum storage   latest
    package contains expected files
    HelmPush.return code should be   0
    ChartMuseum.return code should be   0
    clear chartmuseum storage

    # Repo URL
    push chart directory to url    latest
    HelmPush.return code should be   0
    package exists in chartmuseum storage   latest
    package contains expected files
    HelmPush.return code should be   0
    ChartMuseum.return code should be   0
    clear chartmuseum storage

Chart package can be pushed to ChartMuseum
    # Repo name
    push chart package
    HelmPush.return code should be   0
    package exists in chartmuseum storage
    package contains expected files
    HelmPush.return code should be   0
    ChartMuseum.return code should be   0
    clear chartmuseum storage

    # Repo URL
    push chart package to url
    HelmPush.return code should be   0
    package exists in chartmuseum storage
    package contains expected files
    HelmPush.return code should be   0
    ChartMuseum.return code should be   0
    clear chartmuseum storage

Chart package can be pushed to ChartMuseum with custom version
    # Repo name
    push chart package  latest
    HelmPush.return code should be   0
    package exists in chartmuseum storage   latest
    package contains expected files
    HelmPush.return code should be   0
    ChartMuseum.return code should be   0
    clear chartmuseum storage

    # Repo URL
    push chart package to url  latest
    HelmPush.return code should be   0
    package exists in chartmuseum storage   latest
    package contains expected files
    HelmPush.return code should be   0
    ChartMuseum.return code should be   0
    clear chartmuseum storage

Suite Setup
    set helm version    3
    remove helm plugin
    set helm version    2
    remove helm plugin
    remove chartmuseum logs
    start chartmuseum
    Sleep  2
    set helm version    2
    add chart repo
    set helm version    3
    add chart repo

Suite Teardown
    set helm version    3
    remove helm plugin
    set helm version    2
    remove helm plugin
    remove chart repo
    stop chartmuseum
    print chartmuseum logs
