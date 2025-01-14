FROM golang:1.22.1-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /sample-saas-product

EXPOSE 8080

CMD ["/sample-saas-product"]