import common


class HelmPush(common.CommandRunner):
    def _testchart_path(self):
        if common.USE_OPPOSITE_VERSION:
            if 'helm3' in common.HELM_EXE:
                return '%s/helm2/mychart' % common.TESTCHARTS_DIR
            return '%s/helm3/my-v3-chart' % common.TESTCHARTS_DIR
        else:
            if 'helm3' in common.HELM_EXE:
                return '%s/helm3/my-v3-chart' % common.TESTCHARTS_DIR
            return '%s/helm2/mychart' % common.TESTCHARTS_DIR

    def helm_major_version_detected_by_plugin_is(self, version):
        cmd = '%s push --check-helm-version' % common.HELM_EXE
        self.run_command(cmd)
        self.output_contains(version)

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
