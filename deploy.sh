#!/bin/bash
if [ "$BRANCH" == 'develop' ]; then
    GIT_COMMIT=`echo $COMMIT | head -c 5` 
    scp -r $GOROOT/foundation ottemo@$REMOTE_HOST:~/deploy/foundation-$GIT_COMMIT
    scp $HOME/src/foundation-updater.sh ottemo@REMOTE_HOST:~/deploy/foundation-updater.sh
    ssh ottemo@$REMOTE_HOST "cd /home/ottemo/deploy/ && ln -sf foundation-$GIT_COMMIT foundation-latest" 
fi
