#!/usr/bin/env bash

# vars for script
PWD=$PWD
SRCDIR=/home/ottemo/code/go/src/github.com/ottemo/commerce

cd $SRCDIR
# Build commerce
echo "BUILDING COMMERCE"
./bin/make.sh -tags mongo

# stop the commerce service
echo "STOP COMMERCE SERVICE"
while sudo service ottemo stop >/dev/null 2>&1; do
    echo "warning: commerce is still running"
done
echo "info: commerce has terminated"

echo "DEPLOYING COMMERCE"
# Backup binaries and restart commerce
cp ~/commerce/commerce ~/commerce/commerce.bak
cp commerce ~/commerce/
sudo service ottemo start

# restore PWD
cd $PWD

echo "DEPLOY DONE"
