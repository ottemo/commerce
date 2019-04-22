#!/bin/bash

# build and push commerce image to registry

for i in "$@"
do
case $i in
    version=*)
    version="${i#*=}"
    shift
    ;;
    indocker=*)
    indocker="${i#*=}"
    shift
    ;;
    *)
            # unknown option
    ;;
esac
done

indocker="${indocker:-false}"

GOIMAGE="ottemo/golang:1.7.5" # images that used to build commerce binary
COMMERCEIMAGE="ottemo/commerce"

MYDIR=$(cd `dirname ${BASH_SOURCE[0]}` && pwd)
COMMERCEREPO="$MYDIR/.."
cd $COMMERCEREPO

BUILD=$(git rev-list origin/develop --count)

if ! [ -n "$version" ] ; then
  date=$(date +%Y%m%d-%H%M%S)
  IMAGE="${COMMERCEIMAGE}:${date}-${BUILD}"
else
  IMAGE="${COMMERCEIMAGE}:${version}"
fi
echo "use $IMAGE as image name"

if [ "$indocker" = "true" ] ; then
  echo "build image under docker container"
  echo "generate temporary Dockerfile"
  echo "FROM $GOIMAGE" >Dockerfile.temporary
  echo 'COPY . /go/src/github.com/ottemo/commerce' >>Dockerfile.temporary
  echo 'RUN cd /go/src/github.com/ottemo/commerce && bin/make.sh -tags mongo,redis' >>Dockerfile.temporary

  echo "build commerce in temporary docker container"
  docker build --file Dockerfile.temporary -t temp-commerce-golang .
  if [ $? -ne 0 ]; then
    echo "error in build commerce in temporary docker container"
    exit 2
  fi
  ID=$(docker run -d temp-commerce-golang)
  docker cp $ID:/go/src/github.com/ottemo/commerce/commerce .
else

  echo "build commerce executable with $GOIMAGE docker image"
  docker run -v "$COMMERCEREPO":/go/src/github.com/ottemo/commerce -w /go/src/github.com/ottemo/commerce -e GOOS=linux -e CGO_ENABLED=0 $GOIMAGE bin/make.sh -tags mongo,redis
  if [ $? -ne 0 ]; then
    echo "error in build commerce executable"
    exit 2
  fi
fi

echo "build alpine based commerce container"
docker build -t $IMAGE -t ${COMMERCEIMAGE}:latest .
if [ $? -ne 0 ]; then
  echo "error in build commerce alpine based container"
  exit 2
fi

docker push $IMAGE
if [ $? -ne 0 ]; then
  echo "error in push image"
  exit 2
fi

docker push ${COMMERCEIMAGE}:latest
if [ $? -ne 0 ]; then
  echo "error in push latest commerce image tag"
  exit 2
fi
