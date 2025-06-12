FROM golang:1.16.3-alpine3.13

WORKDIR /app

COPY . .

# Download and install the dependencies:
RUN go get -d -v ./...

# Build the app:
RUN go build -o api .

EXPOSE 8000

CMD ["./api"]