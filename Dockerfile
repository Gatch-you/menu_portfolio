FROM golang:1.20

WORKDIR /app
COPY go.mod .
COPY go.sum .

COPY . .

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

CMD ["air"]