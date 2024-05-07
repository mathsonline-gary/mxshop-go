FROM golang:1.22.2-bookworm

WORKDIR /app

RUN apt-get update
RUN apt-get install -y protobuf-compiler

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# Copy go.mod and go.sum files
COPY go.mod .
COPY go.sum .

# Download Go module dependencies
RUN go mod download

CMD ["tail", "-f", "/dev/null"]