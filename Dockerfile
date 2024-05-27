FROM golang:1.22-alpine

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build

EXPOSE 443
EXPOSE 80

CMD ["./secureforward-proxy"]