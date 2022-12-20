# Stage 1: Build the Go application
FROM golang:1.18-buster as builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY *.go ./
RUN go build -o main .

# Stage 2: Build the final image
FROM gcr.io/distroless/base-debian10

WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 3000
USER nonroot:nonroot

ENTRYPOINT ["./main"]
