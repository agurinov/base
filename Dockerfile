FROM golang:alpine

RUN mkdir -p /bmp/src

ADD . /bmp/src

WORKDIR /bmp/src

RUN set -ex \
		&& go-wrapper download \
		&& go-wrapper install

COPY docker-entrypoint.sh /usr/local/bin/

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["run"]
