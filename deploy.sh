#!/bin/bash
if [ "$BRANCH" == 'develop' ]; then
    GIT_COMMIT=`echo $COMMIT | head -c 5` 
    scp -r $GOPATH/bin/foundation/ ottemo@$REMOTE_HOST:~/deploy/foundation-$GIT_COMMIT
    ssh ottemo@$REMOTE_HOST "cd /home/ottemo/deploy/ && ln -sf foundation-$GIT_COMMIT foundation-latest" 
    ssh ottemo@$REMOTE_HOST "sudo /etc/init.d/ottemo stop && cp ~/foundation/foundation ~/foundation/backup/foundation-$(date +%Y%m%d)"
    ssh ottemo@$REMOTE_HOST "cp ~/deploy/foundation-latest ~/foundation/foundation && sudo /etc/init.d/ottemo start"
fi
