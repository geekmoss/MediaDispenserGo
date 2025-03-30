FROM golang:1.23.2-alpine AS compile

# Set the Current Working Directory inside the container
WORKDIR /app

RUN apk add build-base

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies are cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
ENV CGO_ENABLED=1
RUN go build -o app .

# Stage 2 - Run the app
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the executable from the previous stage
COPY --from=compile /app/app .
ADD config.yaml .

# Expose the application port
EXPOSE 8080

# Command to run the executable
CMD ["./app"]
