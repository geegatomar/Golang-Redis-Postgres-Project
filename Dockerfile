FROM golang:latest

WORKDIR /app 

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

CMD ["go", "run", "main.go"]