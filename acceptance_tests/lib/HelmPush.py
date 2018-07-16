import common


class HelmPush(common.CommandRunner):
    def push_chart_directory(self, version=''):
        cmd = 'helmpush %s/mychart %s' % (common.TESTCHARTS_DIR, common.HELM_REPO_NAME)
        if version:
            cmd += " --version=\"%s\"" % version
        self.run_command(cmd)

    def push_chart_directory_to_url(self, version=''):
        cmd = 'helmpush %s/mychart %s' % (common.TESTCHARTS_DIR, common.HELM_REPO_URL)
        if version:
            cmd += " --version=\"%s\"" % version
        self.run_command(cmd)

    def push_chart_package(self, version=''):
        cmd = 'helmpush %s/mychart/*.tgz %s' % (common.TESTCHARTS_DIR, common.HELM_REPO_NAME)
        if version:
            cmd += " --version=\"%s\"" % version
        self.run_command(cmd)

    def push_chart_package_to_url(self, version=''):
        cmd = 'helmpush %s/mychart %s' % (common.TESTCHARTS_DIR, common.HELM_REPO_URL)
        if version:
            cmd += " --version=\"%s\"" % version
        self.run_command(cmd)
