FROM golang:1.15-alpine AS compile-image
WORKDIR /go/src/mx-api-trainee
COPY go.* /
RUN go mod download
COPY . .
ENV GOOS=linux GOARCH=amd64
# could also set CGO_ENABLED=0, but rumours tells it makes build slower
RUN  go build -ldflags="-s -w" -o /go/bin/mx-api-trainee

FROM alpine:3.13 AS app
# need psql for pinging until db is ready
RUN apk update && apk add postgresql-client
COPY wait-for-postgres.sh /go/src/wait-for-postgres.sh
COPY --from=compile-image /go/bin/mx-api-trainee /go/bin/mx-api-trainee
RUN chmod +x /go/src/wait-for-postgres.sh
ENTRYPOINT  [ "/go/src/wait-for-postgres.sh" ]
# cmd is redundant, running app from entrypoint
# CMD [ "./go/bin/mx-api-trainee" ]