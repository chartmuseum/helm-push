import common


class HelmPush(common.CommandRunner):
    def _testchart_path(self):
        if common.HELM_EXE == 'helm3':
            return '%s/helm3/my-v3-chart' % common.TESTCHARTS_DIR
        return '%s/helm2/mychart' % common.TESTCHARTS_DIR

    def push_chart_directory(self, version=''):
        cmd = '%s push %s %s' % (common.HELM_EXE, self._testchart_path(), common.HELM_REPO_NAME)
        if version:
            cmd += " --version=\"%s\"" % version
        self.run_command(cmd)

    def push_chart_directory_to_url(self, version=''):
        cmd = '%s push %s %s' % (common.HELM_EXE, self._testchart_path(), common.HELM_REPO_URL)
        if version:
            cmd += " --version=\"%s\"" % version
        self.run_command(cmd)

    def push_chart_package(self, version=''):
        cmd = '%s push %s/*.tgz %s' % (common.HELM_EXE, self._testchart_path(), common.HELM_REPO_NAME)
        if version:
            cmd += " --version=\"%s\"" % version
        self.run_command(cmd)

    def push_chart_package_to_url(self, version=''):
        cmd = '%s push %s %s' % (common.HELM_EXE, self._testchart_path(), common.HELM_REPO_URL)
        if version:
            cmd += " --version=\"%s\"" % version
        self.run_command(cmd)
