#!/usr/bin/env bash

WORKDIR=`pwd`
OTTEMODIR="$(cd "$(dirname "$0")" && pwd)"
OTTEMOPKG="github.com/ottemo/commerce"

# select the right version of awk
#    on MacOS and Alpine you may need to install gawk
#
if [[ "$OSTYPE" == "darwin"*  ]] || [ -f /etc/alpine-release ] || [ -f /.dockerenv ] ; then
    AWK=gawk
else
    AWK=awk
fi

# set up env vars for build status
#
DATE=`${AWK} 'BEGIN{ print gensub(/(..)$/, ":\\\1", "g", strftime("%Y-%m-%dT%H:%M:%S%z")); exit }'`
TAGS=""
BUILD=`git -C $OTTEMODIR rev-list origin/develop --count`
BRANCH=`git -C $OTTEMODIR rev-parse --abbrev-ref HEAD`
HASH=`git -C $OTTEMODIR rev-parse --short --verify HEAD`

GOVERSION=`go version | ${AWK} '{print $3}' | awk '{sub(/go/,""); print}'`

# handle flags
#
while test $# -gt 0; do
        case "$1" in
            -tags)
                shift
                TAGS+=$1
                shift
                ;;
            -wd|--work-dir)
                shift
                WORKDIR=$1
                shift
                ;;
            -od|--ottemo-dir)
                shift
                WORKDIR=$1
                shift
                ;;
            *)
                echo $1
                shift
                ;;
        esac
done

cd $WORKDIR

LDFLAGS="-ldflags '"
LDFLAGS+="-X \"github.com/ottemo/commerce/app.buildDate=$DATE\" "
LDFLAGS+="-X \"github.com/ottemo/commerce/app.buildTags=$TAGS\" "
LDFLAGS+="-X \"github.com/ottemo/commerce/app.buildNumber=$BUILD\" "
LDFLAGS+="-X \"github.com/ottemo/commerce/app.buildBranch=$BRANCH\" "
LDFLAGS+="-X \"github.com/ottemo/commerce/app.buildHash=$HASH\" "
LDFLAGS+="'"

# uncomment for go versions previous to 1.7 using the old linker
#
#  LDFLAGS=${LDFLAGS//=/\" \"}

if [ -z "$GOPATH" ]; then
REPLACE="/src/$OTTEMOPKG"
export GOPATH="${OTTEMODIR/$REPLACE/}"
fi


if [ -n "$TAGS" ]; then
TAGS=$(echo $TAGS| sed 's/,/ /g')
TAGS="-tags '$TAGS'"
fi

# install glide to $GOPATH
#     if it does not exist yet
#
#if [ ! -f "$GOPATH/bin/glide" ]; then
#    echo "Glide not found, installing it.\n"
#    go get github.com/Masterminds/glide
#fi

# install project dependencies
$GOPATH/bin/glide install


CMD="go build -a $TAGS $LDFLAGS $OTTEMOPKG"
eval CGO_ENABLED=0 $CMD
