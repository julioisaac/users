FROM golang:1.20.6-alpine3.18 AS builder

WORKDIR /users

COPY . .

RUN go mod download && go mod verify

RUN apk add --no-cache git

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags='-w -s -extldflags "-static"' -v -a -o /go/bin/users .

FROM alpine:3.18

COPY --from=builder /users/docs /docs
COPY --from=builder /go/bin/users /go/bin/users

EXPOSE 9000

RUN adduser -u 1000 -D users-app
RUN chown -R 1000  /go/bin/users

USER users-app

ENTRYPOINT ["/go/bin/users"]