FROM abhishekkr/alpine:dev
MAINTAINER AbhishekKr <abhikumar163@gmail.com>

ENV DORY_BASEDIR /opt/dory

RUN mkdir -p $DORY_BASEDIR

COPY ./bin/dory-linux-amd64 $DORY_BASEDIR/dory
COPY ./templates $DORY_BASEDIR/templates
COPY ./w3assets $DORY_BASEDIR/w3assets

EXPOSE 8080

WORKDIR $DORY_BASEDIR

ENTRYPOINT ./dory

