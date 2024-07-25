FROM golang:1.22.4-alpine3.19 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY .. .
RUN go build -o ./app ./cmd/api/.

FROM alpine:3.19 AS publish

WORKDIR /app

COPY --from=build ./app .

EXPOSE 8080

ENTRYPOINT [ "./app" ]