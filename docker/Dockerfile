FROM ubuntu:latest

RUN apt-get update && apt-get install -y \
build-essential \
curl \
make

# Install Go
ARG GO_VERSION=1.20
RUN curl -LO "https://golang.org/dl/go$GO_VERSION.linux-amd64.tar.gz" \
    && tar -C /usr/local -xzf "go$GO_VERSION.linux-amd64.tar.gz" \
    && rm "go$GO_VERSION.linux-amd64.tar.gz"

# Set environment variables for Go
ENV GOPATH="/go"
ENV PATH="/usr/local/go/bin:${GOPATH}/bin:${PATH}"

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the Go source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port on which the server will listen
EXPOSE 8080

# Set the command to run the Go application when the container starts
CMD ["./main"]
