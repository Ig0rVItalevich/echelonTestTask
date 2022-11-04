FROM golang:1.18-buster

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o echelon-test ./cmd/server/main.go
RUN go build -o client ./cmd/client/main.go

CMD ["echelon-test"]