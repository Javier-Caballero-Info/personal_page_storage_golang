FROM golang:1.8

WORKDIR /go/src/app

ENV PORT=3000

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 3000

CMD ["app"]