FROM golang:1.15-alpine AS compile-image
WORKDIR /go/src/mx-api-trainee
COPY go.* /
RUN go mod download
COPY . .
ENV GOOS=linux GOARCH=amd64
RUN  go build -ldflags="-s -w" -o /go/bin/mx-api-trainee

FROM alpine:3.13 AS app
COPY --from=compile-image /go/bin/mx-api-trainee /go/bin/mx-api-trainee
ENTRYPOINT  [ "/go/bin/mx-api-trainee" ]