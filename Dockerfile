FROM golang:1.11-alpine AS build-env

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH   

# now copy your app to the proper build path
RUN mkdir -p $GOPATH/src/github.com/asapp-noclist
ADD . $GOPATH/src/github.com/asapp-noclist

# should be able to build now
WORKDIR $GOPATH/src/github.com/asapp-noclist
RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/asapp-noclist/asapp-noclist .
ENTRYPOINT ["./asapp-noclist"]