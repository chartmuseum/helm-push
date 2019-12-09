import os

import common

class Helm(common.CommandRunner):
    def set_helm_version(self, version):
        version = str(version)
        if version == '2':
            common.HELM_EXE = 'HELM_HOME=%s helm2' % os.getenv('TEST_HELM_HOME', '')
        elif version == '3':
            common.HELM_EXE = 'helm3'
        else:
            raise Exception('invalid Helm version provided: %s' % version)

    def add_chart_repo(self):
        self.remove_chart_repo()
        self.run_command('%s repo add %s %s' % (common.HELM_EXE, common.HELM_REPO_NAME, common.HELM_REPO_URL))

    def remove_chart_repo(self):
        self.run_command('%s repo remove %s' % (common.HELM_EXE, common.HELM_REPO_NAME))

    def install_helm_plugin(self):
        self.run_command('%s plugin install %s' % (common.HELM_EXE, self.rootdir))

    def check_helm_plugin(self):
        self.run_command('%s plugin list | grep ^push' % common.HELM_EXE)

    def run_helm_plugin(self):
        self.run_command('%s push --help' % common.HELM_EXE)

    def remove_helm_plugin(self):
        self.run_command('%s plugin remove push' % common.HELM_EXE)
