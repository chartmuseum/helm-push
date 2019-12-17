import os
import shutil

import common


class ChartMuseum(common.CommandRunner):
    def start_chartmuseum(self):
        self.stop_chartmuseum()
        os.chdir(self.rootdir)
        shutil.rmtree(common.STORAGE_DIR, ignore_errors=True)
        cmd = 'chartmuseum --debug --port=%d --storage="local" ' % common.PORT
        cmd += '--storage-local-rootdir=%s >> %s 2>&1' % (common.STORAGE_DIR, common.LOGFILE)
        print(cmd)
        self.run_command(cmd, detach=True)

    def stop_chartmuseum(self):
        self.run_command('pkill -9 chartmuseum')
        shutil.rmtree(common.STORAGE_DIR, ignore_errors=True)

    def remove_chartmuseum_logs(self):
        os.chdir(self.rootdir)
        self.run_command('rm -f %s' % common.LOGFILE)

    def print_chartmuseum_logs(self):
        os.chdir(self.rootdir)
        self.run_command('cat %s' % common.LOGFILE)

    def package_exists_in_chartmuseum_storage(self, find=''):
        self.run_command('find %s -maxdepth 1 -name "*%s.tgz" | grep tgz' % (common.STORAGE_DIR, find))

    def package_contains_expected_files(self):
        if 'helm3' in common.HELM_EXE:
            # Check for values.schema.json in Helm 3 (a Helm 3-specific file)
            self.run_command('(cd %s && mkdir -p tmp && tar -xf *.tgz --directory tmp && find tmp -name values.schema.json | grep values.schema.json)' % common.STORAGE_DIR)
        else:
            # Check for requirements.yaml in Helm 2 (a Helm 2-specific file)
            self.run_command('(cd %s && mkdir -p tmp && tar -xf *.tgz --directory tmp && find tmp -name grep requirements.yaml | grep requirements.yaml)' % common.STORAGE_DIR)

    def clear_chartmuseum_storage(self):
        self.run_command('rm %s*.tgz' % common.STORAGE_DIR)
