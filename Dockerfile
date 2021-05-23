FROM golang:alpine AS build

# Add Maintainer Info
LABEL maintainer="guillaume.penaud@gmail.com"

RUN \
  apk add --no-cache git &&\
  mkdir /application

ADD . /application/
WORKDIR /application

# Download all the dependencies
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux \
  go build \
    -a -installsuffix cgo \
    -o /needys-api-need \
  /application/cmd/needys-api-need-server/main.go

# ---------------------------------------------------------------------------- #

FROM alpine:latest

RUN adduser --system --disabled-password --home /needys-api-need needys-api-need

WORKDIR /needys-api-need
USER needys-api-need

COPY --from=build /needys-api-need .

EXPOSE 8010

CMD ["./needys-api-need"]
