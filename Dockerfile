FROM golang

RUN mkdir -p /bmp/src

ADD . /bmp/src

WORKDIR /bmp/src

RUN set -ex \
		&& go get -v -d ./... \
		&& go install -v ./...

COPY docker-entrypoint.sh /usr/local/bin/

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["run"]
