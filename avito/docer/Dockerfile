FROM golang:lasted

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download \
RUM go build - app ./cmd/app.go

CMD ["./app"]