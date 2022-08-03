FROM golang:alpine as builder

WORKDIR /app

RUN ls -l

ADD . .

RUN go mod download

RUN go build -o bin/api-service ./cmd/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/bin /app

CMD [ "./api-service" ]
