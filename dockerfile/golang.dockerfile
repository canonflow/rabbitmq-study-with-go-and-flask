FROM golang:1.23.3

WORKDIR /app

# Install FFmpeg and clean up in one layer to reduce image size
RUN apt-get update && \
    apt-get install -y --no-install-recommends dos2unix && \
    rm -rf /var/lib/apt/lists/*

# Copy go mod/sum first for caching
COPY go/go.mod go/go.sum ./

RUN go mod download

# Copy all source code
COPY go ./

# Build both binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-service ./main.go

# Copy startup script
COPY scripts/go-start.sh ./
RUN chmod +x ./go-start.sh
# RUN dos2unix go-start.sh

EXPOSE 8990

# Run both services via startup script
CMD ["./go-start.sh"]