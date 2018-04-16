import common


class Helm(common.CommandRunner):
    def add_chart_repo(self):
        self.remove_chart_repo()
        self.run_command('helm repo add %s %s' % (common.HELM_REPO_NAME, common.HELM_REPO_URL))

    def remove_chart_repo(self):
        self.run_command('helm repo remove %s' % common.HELM_REPO_NAME)

    def install_helm_plugin(self):
        self.run_command('helm plugin install %s' % self.rootdir)

    def check_helm_plugin(self):
        self.run_command('helm plugin list | grep ^push')

    def run_helm_plugin(self):
        self.run_command('helm push --help')

    def remove_helm_plugin(self):
        self.run_command('helm plugin remove push')
