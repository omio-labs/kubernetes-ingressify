FROM golang:1.9

WORKDIR /go/src/github.com/goeuro/ingress-generator-kit
COPY Godeps Godeps
COPY bin/init.sh bin/
RUN bin/init.sh

COPY . /go/src/github.com/goeuro/ingress-generator-kit
CMD ["bin/test.sh"]
