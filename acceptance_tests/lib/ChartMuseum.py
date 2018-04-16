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

    def clear_chartmuseum_storage(self):
        self.run_command('rm %s*.tgz' % common.STORAGE_DIR)
