# Start from the official golang image
FROM golang:1.22 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
RUN go get -u github.com/cosmtrek/air
# Copy the source code from the current directory to the Working Directory inside the container
COPY . .

# Command to run the executable
CMD ["air"]