FROM golang:1.22-alpine
WORKDIR /tracker
COPY . ./
RUN ls
RUN go mod download -x
RUN go build AuthService/main.go
ENTRYPOINT [ "./main" ]