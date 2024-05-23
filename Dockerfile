# Use the official Golang image
FROM golang:latest

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go application source code to the container
COPY . .

RUN go mod download

# Build the Go application
RUN go build -o soda-go

# Install FFmpeg
RUN apt-get update && apt-get install -y ffmpeg

# Expose the port the API will run on
EXPOSE 8080

# Command to run the Go application
CMD ["./soda-go"]