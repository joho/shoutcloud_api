FROM golang:1.3.1

# Folder is created by golang image, we need it out the say to put
# our own src bin in as volumes
RUN rm -r /go/*

ENV GOBIN /go/bin

ENV PORT 8080

WORKDIR /go/src/github.com/joho/shoutcloud_api

ADD provision.sh /provision.sh

CMD ["/go/bin/shoutcloud_api"]

EXPOSE 8080
