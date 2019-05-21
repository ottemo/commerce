# TODO: update this old bootstrap script for setting up commerce the first time in /opt dir
apt-get -y install bzr
apt-get -y install git
apt-get -y install golang-go


rm -rf /opt/ottemo
mkdir -pv /opt/ottemo/go/src/github.com/ottemo
mkdir -pv /opt/ottemo/go/bin
mkdir -pv /opt/ottemo/go/pkg
mkdir -pv /opt/ottemo/media

export GOPATH=/opt/ottemo/go

git clone https://github.com/ottemo/commerce /opt/ottemo/go/src/github.com/ottemo/commerce

cd $GOPATH/bin
echo "media.fsmedia.folder=/opt/ottemo/media" >> ottemo.ini
echo "mongodb.db=ottemo-demo" >> ottemo.ini
echo "mongodb.uri=mongodb://DB_USER:DB_PASSWROD@MONGO_DB_URI:27017/ottemo" >> ottemo.ini

cd $GOPATH/src/github.com/ottemo/commerce && go get -t
cd $GOPATH/src/github.com/ottemo/commerce && go get gopkg.in/mgo.v2
cd $GOPATH/src/github.com/ottemo/commerce && go get gopkg.in/mgo.v2/bson
cd $GOPATH/src/github.com/ottemo/commerce && go build -tags mongo
