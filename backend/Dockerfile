# ----Build stage----
FROM golang:1.23.0-alpine AS build

WORKDIR /go/services/utc

COPY ./go.mod ./go.sum /go/services/utc/
RUN go mod download && \
    go mod tidy

# Install curl and air.
RUN apk update && \
    apk upgrade && \
    apk add --no-cache curl && \
    go install github.com/air-verse/air@latest

COPY . .

RUN go build -o app cmd/main.go

# ----Execution stage----
FROM alpine:latest

WORKDIR /approot

COPY --from=build /go/services/utc/app ./

RUN addgroup go && \
    adduser -D -G go go && \
    chown -R go:go /approot/app

USER go

CMD /approot/app