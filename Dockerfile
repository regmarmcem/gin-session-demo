FROM golang:1.20-alpine AS builder
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o main

FROM gcr.io/distroless/static-debian11
USER nonroot
WORKDIR /app
COPY --from=builder /go/src/app/main /go/src/app/.env ./
EXPOSE 8080
CMD ["./main"]

