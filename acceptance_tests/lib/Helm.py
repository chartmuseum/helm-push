import os

import common

class Helm(common.CommandRunner):
    def set_helm_version(self, version):
        version = str(version)
        if version == '3':
            common.HELM_EXE = 'XDG_CACHE_HOME=%s XDG_CONFIG_HOME=%s XDG_DATA_HOME=%s helm3' % \
                (os.getenv('TEST_V3_XDG_CACHE_HOME', ''), os.getenv('TEST_V3_XDG_CONFIG_HOME', ''),
                    os.getenv('TEST_V3_XDG_DATA_HOME', ''))
        elif version == '4':
            common.HELM_EXE = 'XDG_CACHE_HOME=%s XDG_CONFIG_HOME=%s XDG_DATA_HOME=%s helm4' % \
                (os.getenv('TEST_V3_XDG_CACHE_HOME', ''), os.getenv('TEST_V3_XDG_CONFIG_HOME', ''),
                    os.getenv('TEST_V3_XDG_DATA_HOME', ''))
        else:
            raise Exception('invalid Helm version provided: %s' % version)

    def use_test_chart_built_by_same_helm_version(self):
        common.USE_OPPOSITE_VERSION = False

    def use_test_chart_built_by_opposite_helm_version(self):
        common.USE_OPPOSITE_VERSION = True

    def add_chart_repo(self):
        self.remove_chart_repo()
        self.run_command('%s repo add %s %s' % (common.HELM_EXE, common.HELM_REPO_NAME, common.HELM_REPO_URL))

    def remove_chart_repo(self):
        self.run_command('%s repo remove %s' % (common.HELM_EXE, common.HELM_REPO_NAME))

    def setup_plugin_manifest(self):
        # In development mode, manually copy the correct manifest for this Helm version
        if 'helm4' in common.HELM_EXE:
            self.run_command('cp %s/testdata/plugin-helm4.yaml %s/plugin.yaml' % (self.rootdir, self.rootdir))
        elif 'helm3' in common.HELM_EXE:
            self.run_command('cp %s/testdata/plugin-helm3.yaml %s/plugin.yaml' % (self.rootdir, self.rootdir))

    def install_helm_plugin(self):
        # Setup correct manifest before installing
        self.setup_plugin_manifest()

        # Set HELM_MAJOR_VERSION env var to help install script detect correct version
        if 'helm4' in common.HELM_EXE:
            helm_version = '4'
        elif 'helm3' in common.HELM_EXE:
            helm_version = '3'
        else:
            helm_version = '3'  # Default to 3
        self.run_command('HELM_MAJOR_VERSION=%s %s plugin install %s' % (helm_version, common.HELM_EXE, self.rootdir))

    def check_helm_plugin(self):
        self.run_command('%s plugin list | grep ^cm-push' % common.HELM_EXE)

    def run_helm_plugin(self):
        self.run_command('%s cm-push --check-helm-version' % common.HELM_EXE)

    def remove_helm_plugin(self):
        self.run_command('%s plugin remove cm-push' % common.HELM_EXE)
