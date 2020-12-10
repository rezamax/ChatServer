# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from golang:1.12-alpine base image
FROM golang:1.15-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

ENV GOPATH=/app

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh
RUN go get github.com/gorilla/websocket
RUN go get -u github.com/gorilla/mux



# Copy the source from the current directory to the Working Directory inside the container
COPY . .

#RUN go run main