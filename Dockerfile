FROM golang

RUN mkdir -p /go/src/app

ADD . /go/src/app

WORKDIR /go/src/app

RUN set -ex \
		&& go get -v -d ./... \
		&& go install -v ./... \
		\
		&& mv /go/src/app/docker-entrypoint.sh /usr/local/bin \
		&& chmod +x /usr/local/bin/docker-entrypoint.sh

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["run"]
