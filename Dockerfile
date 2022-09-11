FROM golang:1.18-alpine as gobuild
WORKDIR /src
COPY . .
RUN go get -d -v
RUN go build -ldflags "-s -w" -o csdb

FROM alpine:latest
COPY --from=gobuild /src/csdb /csdb

#Map volume as /opt for app-data
ENTRYPOINT ["/csdb"]

