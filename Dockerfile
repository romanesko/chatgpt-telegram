FROM golang:latest

# Copy the application source code into the container
COPY app /app

# Set the working directory
WORKDIR /app

# Install any dependencies needed by the application
RUN go get -d -v ./...
RUN go install -v ./...

# Build the application
RUN go build -o app

# Expose the port that the application listens on
#EXPOSE 8080

# Start the application
CMD ["./app"]