## We specify the base image we need for our
## go application
##FROM golang:1.12.0-alpine3.9

FROM golang:1.16-alpine

## We create an /app directory within our
## image that will hold our application source
## files
RUN mkdir /app

## We copy everything in the root directory
## into our /app directory
##ADD . /app
COPY go.mod /app
COPY go.sum /app
COPY webserver.go /app
COPY post.html /app
COPY create_post.html /app

## We specify that we now wish to execute 
## any further commands inside our /app
## directory
WORKDIR /app

RUN go mod download

## we run go build to compile the binary
## executable of our Go program
RUN go build -o webserver .

EXPOSE 8080

## Our start command which kicks off
## our newly created binary executable
CMD ["/app/webserver"]

