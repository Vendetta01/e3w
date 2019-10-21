ARG E3W_GIT_SRC_URL=https://github.com/VendettA01/e3w.git
#ARG E3W_GIT_COMMIT=c85c4e78f43761cc070c082158ca2c5b6e895ed9


FROM golang:1.9 AS backend

ARG E3W_GIT_SRC_URL
#ENV E3W_GIT_COMMIT

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN mkdir -p /go/src/github.com/VendettA01/e3w
RUN git clone ${E3W_GIT_SRC_URL} /go/src/github.com/VendettA01/e3w
WORKDIR /go/src/github.com/VendettA01/e3w
#RUN git checkout ${E3W_GIT_COMMIT}
RUN dep ensure
RUN CGO_ENABLED=0 go build

RUN git clone https://github.com/vishnubob/wait-for-it.git /tmp/wait-for-it/


FROM node:8 AS frontend

ARG E3W_GIT_SRC_URL
#ENV E3W_GIT_COMMIT

RUN mkdir /app
RUN git clone ${E3W_GIT_SRC_URL} /tmp/e3w
WORKDIR /tmp/e3w
#RUN git checkout ${E3W_GIT_COMMIT}
RUN cp -a /tmp/e3w/static/* /app/
#ADD static /app
WORKDIR /app
RUN npm --registry=https://registry.npm.taobao.org \
--cache=$HOME/.npm/.cache/cnpm \
--disturl=https://npm.taobao.org/mirrors/node \
--userconfig=$HOME/.cnpmrc install && npm run publish

FROM alpine:edge

RUN apk --no-cache add bash
RUN mkdir -p /app/static/dist /app/conf
COPY --from=backend /go/src/github.com/VendettA01/e3w/e3w /app
COPY --from=frontend /app/dist /app/static/dist
COPY --from=backend /go/src/github.com/VendettA01/e3w/conf/config.default.ini /app/conf
COPY --from=backend /tmp/wait-for-it/wait-for-it.sh /usr/bin/wait-for-it.sh
RUN chmod 755 /usr/bin/wait-for-it.sh
COPY scripts/* /usr/bin/
EXPOSE 8080
WORKDIR /app

ENTRYPOINT ["docker_entrypoint.sh"]

