FROM dockerfile/go

MAINTAINER n4sjamk

RUN apt-get update && \
	apt-get install -y mercurial

ADD . /gopath/src/teamboard-crypt

RUN cd /gopath/src/teamboard-crypt && \
	go get

CMD teamboard-crypt
