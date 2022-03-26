FROM golang:1.18

WORKDIR /go/src/github.com/goeuro/kubernetes-ingressify
COPY Godeps Godeps
COPY bin/init.sh bin/
RUN bin/init.sh

COPY . /go/src/github.com/goeuro/kubernetes-ingressify
CMD ["bin/test.sh"]
