# FROM golang:latest 

# RUN mkdir /go_service 

# WORKDIR /go_service
# COPY go_service.go .

# RUN go mod init github.com/andrew/go_service
# RUN go get github.com/gorilla/mux
# RUN go get github.com/google/uuid
# RUN go get github.com/go-redis/redis/v8

# RUN go get github.com/adjust/rmq
# RUN go build -o main .
# CMD ["go run /app/main.go"]


# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add ca-certificates bash

WORKDIR /root/

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.3/wait /wait
RUN chmod +x /wait

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8000 to the outside world
EXPOSE 8000

# Command to run the executable
CMD ["bash", "-c", "./wait && ./main"] 