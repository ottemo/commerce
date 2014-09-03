FROM ottemo/golang:v1

# Add SSH keys
#ADD $HOME/.ssh/id_rsa /root/.ssh/id_rsa
#RUN echo "IdentityFile ~/.ssh/id_rsa" >> /etc/ssh/ssh_config
#RUN mkdir -p /root/.ssh
#ADD id_rsa ~/.ssh/id_rsa
#RUN ssh-keyscan github.com >> ~/.ssh/known_hosts

# Build Foundation
RUN mkdir -pv /root/go/{bin,src/{github.com/ottemo},pkg}
RUN git clone https://ottemo-dev:freshbox111222333@github.com/ottemo/foundation.git /root/go/src/github.com/ottemo/foundation

RUN cd $GOPATH/src/github.com/ottemo/foundation && go get -t 
RUN cd $GOPATH/src/github.com/ottemo/foundation && go build && go install
