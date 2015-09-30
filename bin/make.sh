#!/bin/bash

WORKDIR=`pwd`
OTTEMODIR="$(cd "$(dirname "$0")" && pwd)"
OTTEMOPKG="github.com/ottemo/foundation"
if [[ "$OSTYPE" == "darwin"*  ]]; then
    DATE=`gawk 'BEGIN{ print gensub(/(..)$/, ":\\\1", "g", strftime("%Y-%m-%dT%H:%M:%S%z")); exit }'`
else
    DATE=`awk 'BEGIN{ print gensub(/(..)$/, ":\\\1", "g", strftime("%Y-%m-%dT%H:%M:%S%z")); exit }'`
fi
TAGS=""
BUILD=`git -C $OTTEMODIR rev-list origin/develop --count`
BRANCH=`git -C $OTTEMODIR rev-parse --abbrev-ref HEAD`
HASH=`git -C $OTTEMODIR rev-parse --short --verify HEAD`

if [[ "$OSTYPE" == "darwin"*  ]]; then
    GOVERSION=`go version | gawk '{print gensub(/.*go([0-9]+[.][0-9]+).[0-9]+.*/, "\\\1", "1")}'`
else
    GOVERSION=`go version | awk '{print gensub(/.*go([0-9]+[.][0-9]+).[0-9]+.*/, "\\\1", "1")}'`
fi

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
LDFLAGS+="-X \"github.com/ottemo/foundation/app.buildDate=$DATE\" "
LDFLAGS+="-X \"github.com/ottemo/foundation/app.buildTags=$TAGS\" "
LDFLAGS+="-X \"github.com/ottemo/foundation/app.buildNumber=$BUILD\" "
LDFLAGS+="-X \"github.com/ottemo/foundation/app.buildBranch=$BRANCH\" "
LDFLAGS+="-X \"github.com/ottemo/foundation/app.buildHash=$HASH\" "
LDFLAGS+="'"

if [ "`awk "BEGIN{ if ($GOVERSION < 1.5) print 1 }"`" == "1" ]; then
  LDFLAGS=${LDFLAGS//=/\" \"}
fi

if [ -z "$GOPATH" ]; then 
REPLACE="/src/$OTTEMOPKG"
export GOPATH="${OTTEMODIR/$REPLACE/}"
fi


if [ -n "$TAGS" ]; then
TAGS="-tags $TAGS"
fi

CMD="go get $TAGS $OTTEMOPKG"
eval $CMD

CMD="go build -a $TAGS $LDFLAGS $OTTEMOPKG"
eval $CMD
