FROM golang:1.14 as builder
ADD . /go/src/github.com/digitalocean/artifactory-resource
WORKDIR /go/src/github.com/digitalocean/artifactory-resource
RUN make build

FROM alpine:3.11 as resource
RUN apk add --update --no-cache bash bash-completion openssh git
COPY --from=builder /go/src/github.com/digitalocean/artifactory-resource/build /opt/resource
RUN ln -s /opt/resource/get /opt/resource/in && ln -s /opt/resource/put /opt/resource/out && chmod +x /opt/resource/*
CMD ["/bin/bash"]

FROM resource
LABEL MAINTAINER=digitalocean
