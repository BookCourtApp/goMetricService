# Use the official Go image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code into the container
COPY . .

# Build the Go application
RUN go build -o ./cmd/goMetricService ./cmd/goMetricService.go 

# Specify the command to run your application
CMD ["./cmd/goMetricService","-config=./config/local.yaml"]