FROM golang:alpine AS builder

WORKDIR /app

# Copy Go source code
COPY . .

# Install dependencies
RUN go mod download

# Build the Go binary
RUN go build -o stocktrend .

# Create a slim runtime image
FROM alpine AS runtime

WORKDIR /app

# Copy the built binary
COPY --from=builder /app/stocktrend /app/stocktrend

# Set the entry point
CMD ["/app/stocktrend"]
