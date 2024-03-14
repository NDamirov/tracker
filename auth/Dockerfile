FROM golang:1.22-alpine
WORKDIR /auth
COPY go.* ./
RUN go mod download -x
COPY * ./
RUN go build .
ENTRYPOINT [ "./main" ]