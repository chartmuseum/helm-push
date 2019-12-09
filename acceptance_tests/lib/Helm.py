import common

HELM_EXE = 'helm2'

class Helm(common.CommandRunner):
    def set_helm_version(self, version):
        global HELM_EXE
        version = str(version)
        if version == '2':
            HELM_EXE = 'helm2'
        elif version == '3':
            HELM_EXE = 'helm3'
        else:
            raise Exception('invalid Helm version provided: %s' % version)

    def add_chart_repo(self):
        self.remove_chart_repo()
        self.run_command('%s repo add %s %s' % (HELM_EXE, common.HELM_REPO_NAME, common.HELM_REPO_URL))

    def remove_chart_repo(self):
        self.run_command('%s repo remove %s' % (HELM_EXE, common.HELM_REPO_NAME))

    def install_helm_plugin(self):
        self.run_command('%s plugin install %s' % (HELM_EXE, self.rootdir))

    def check_helm_plugin(self):
        self.run_command('%s plugin list | grep ^push' % HELM_EXE)

    def run_helm_plugin(self):
        self.run_command('%s push --help' % HELM_EXE)

    def remove_helm_plugin(self):
        self.run_command('%s plugin remove push' % HELM_EXE)
