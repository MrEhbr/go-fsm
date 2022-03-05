# build
FROM            golang:1.16-alpine as builder
# dynamic config
ARG             BUILD_DATE
ARG             VCS_REF
ARG             VERSION

RUN             apk add --no-cache git gcc musl-dev make bash
ENV             GO111MODULE=on
WORKDIR         /go/src/github.com/MrEhbr/go-fsm
COPY            go.* ./
RUN             go mod download
COPY            . ./
RUN             make install VCS_REF=$VCS_REF VERSION=$VERSION BUILD_DATE=$BUILD_DATE

# minimalist runtime
FROM alpine:3.11
# dynamic config
ARG             BUILD_DATE
ARG             VCS_REF
ARG             VERSION

LABEL org.label-schema.build-date=$BUILD_DATE \
    org.label-schema.name="go-fsm" \
    org.label-schema.description="" \
    org.label-schema.url="" \
    org.label-schema.vcs-ref=$VCS_REF \
    org.label-schema.vcs-url="https://github.com/MrEhbr/go-fsm" \
    org.label-schema.vendor="Alexey Burmistrov" \
    org.label-schema.version=$VERSION \
    org.label-schema.schema-version="1.0" \
    org.label-schema.cmd="docker run -i -t --rm MrEhbr/go-fsm" \
    org.label-schema.help="docker exec -it $CONTAINER go-fsm --help"
COPY            --from=builder /go/bin/go-fsm /bin/
ENTRYPOINT      ["/bin/go-fsm"]
#CMD             []
