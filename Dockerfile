FROM golang

RUN mkdir -p /go/src/app

ADD . /go/src/app

WORKDIR /go/src/app

RUN set -ex \
		&& go get -v -d ./... \
		&& go install -v ./...

COPY docker-entrypoint.sh /usr/local/bin/

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["run"]
