FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN apk add git
RUN go mod download
ADD . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-one-server ./cmd/go-one-server/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/ /app/
EXPOSE 8081
EXPOSE 8081
CMD ["./go-one-server"]