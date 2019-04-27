FROM golang

ADD . /go/src/github.com/TGRHavoc/gello
WORKDIR /go/src/github.com/TGRHavoc/gello

# Make sure the vendor packages have their dependencies met
RUN go get -t -v ./vendors/... 

RUN go build -v

ENTRYPOINT [ "gello" ]