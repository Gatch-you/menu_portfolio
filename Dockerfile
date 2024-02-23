FROM golang:1.20

WORKDIR /app
COPY go.mod .
COPY go.sum .

COPY . .

RUN go build -buildvcs=false -o /usr/local/bin/myapp .

CMD ["/usr/local/bin/myapp"]