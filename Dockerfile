FROM golang:1.22-alpine

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod .
COPY go.sum .

# Download Go module dependencies
RUN go mod download

# Copy the start-container script
COPY start-container .
COPY start-container /usr/local/bin/start-container

# Make the start-container script executable
RUN chmod +x start-container
RUN chmod +x /usr/local/bin/start-container

# Expose port 8000
EXPOSE 8000

CMD ["tail", "-f", "/dev/null"]