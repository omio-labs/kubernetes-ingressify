FROM golang:1.9

WORKDIR /go/src/github.com/goeuro/ingress-generator-kit
COPY . /go/src/github.com/goeuro/ingress-generator-kit
RUN bin/init.sh
CMD ["bin/test.sh"]
