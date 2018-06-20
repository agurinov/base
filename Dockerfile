FROM golang:alpine

RUN mkdir -p /bmp/src

ADD . /bmp/src

WORKDIR /bmp/src

RUN set -ex \
		&& go get -v -d ./... \
		&& go install -v ./... \
			-o app

COPY docker-entrypoint.sh /usr/local/bin/

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["run"]
