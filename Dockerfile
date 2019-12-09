FROM golang:1.9 AS backend

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
COPY src/ /go/src/github.com/VendettA01/e3w/src
WORKDIR /go/src/github.com/VendettA01/e3w/src
RUN dep ensure
RUN CGO_ENABLED=0 go build


FROM node:8 AS frontend

RUN mkdir /app
COPY static/ /app/
WORKDIR /app
RUN npm --registry=https://registry.npm.taobao.org \
--cache=$HOME/.npm/.cache/cnpm \
--disturl=https://npm.taobao.org/mirrors/node \
--userconfig=$HOME/.cnpmrc install && npm run publish

FROM confd:latest

RUN apk --no-cache add bash supervisor
RUN mkdir -p /app/static/dist /app/conf
COPY --from=backend /go/src/github.com/VendettA01/e3w/src/e3w /app
COPY --from=frontend /app/dist /app/static/dist
COPY --from=backend /go/src/github.com/VendettA01/e3w/src/conf/config.default.ini /app/conf
COPY scripts/* /usr/bin/

COPY etc/ /etc/

EXPOSE 8080
WORKDIR /app

ENTRYPOINT ["docker_entrypoint.sh"]

