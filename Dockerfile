FROM golang:latest
COPY . /app
WORKDIR /app
RUN go build -o rabbitsky-server
CMD ./rabbitsky-server

