import os
from fabric.api import *

env.hosts = ['kullo2.kullo.net']
env.user = 'root'

SERVER_DIR = '/opt/crashreports'
EXE_NAME = 'crashreports'

@task(default=True)
def deploy():
	local('make')
	local('make test')
	with cd(SERVER_DIR):
		put(EXE_NAME, EXE_NAME + '-new', mode=0755)
		run('rm ' + EXE_NAME + '-old', warn_only=True)
		run('[ -e "' + EXE_NAME + '" ] || touch "' + EXE_NAME + '"')
		run('mv ' + EXE_NAME + ' ' + EXE_NAME + '-old && mv ' + EXE_NAME + '-new ' + EXE_NAME)
		run('systemctl stop ' + EXE_NAME, warn_only=True)
		run('systemctl start ' + EXE_NAME)

@task
def back_to_last_version():
	with cd(SERVER_DIR):
		run('mv ' + EXE_NAME + ' ' + EXE_NAME + '-new && mv ' + EXE_NAME + '-old ' + EXE_NAME)
		run('systemctl stop ' + EXE_NAME, warn_only=True)
		run('systemctl start ' + EXE_NAME)

