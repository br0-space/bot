FROM golang:1.15-alpine as build

# Build-time metadata as defined at http://label-schema.org
ARG CI_COMMIT_SHA
ARG CI_COMMIT_REF_SLUG
LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.name="kmptnz telegram bot" \
      org.label-schema.description="Telegram bot written in golang" \
      org.label-schema.url="https://github.com/NeoVG/kmptnzbot" \
      org.label-schema.vendor="NeoVG" \
      org.label-schema.version=$CI_COMMIT_REF_SLUG \
      org.label-schema.schema-version="1.0"

# Build-time req.
RUN apk --no-cache add git

# go config
ENV GO111MODULE=on
ENV CGO_ENABLED=0

WORKDIR /go/src/app
COPY . .

RUN go build -o bin/kmptnzbot main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /opt/kmptnzbot

COPY --from=build /go/src/app/bin/kmptnzbot /opt/kmptnzbot/bin/kmptnzbot

CMD ["/opt/kmptnzbot/bin/kmptnzbot"]
