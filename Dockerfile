FROM google/golang:latest

RUN go get -u github.com/tools/godep

ADD . $GOPATH/src/github.com/AlexisMontagne/proxyd
RUN godep get github.com/AlexisMontagne/proxyd/...
RUN cd $GOPATH/src/github.com/AlexisMontagne/proxyd && godep restore
