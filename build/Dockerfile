FROM golang:1.21-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/shortener ./cmd/server

FROM alpine:3.18
COPY --from=build /bin/shortener /bin/shortener
EXPOSE 8080
ENTRYPOINT ["/bin/shortener"]
