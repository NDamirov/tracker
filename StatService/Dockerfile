FROM golang:1.22-alpine
WORKDIR /tracker
COPY . ./
RUN ls
RUN go mod download -x
RUN go build StatService/main.go
ENTRYPOINT [ "./main" ]