# Stage 1

FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY * ./

RUN export CGO_ENABLED=0 && go build -o beerserver main.go

# Stage 2

FROM alpine:latest AS runner

WORKDIR /

COPY --from=builder /app/beerserver /beerel/beerserver
COPY --from=builder /app/templates/index.html /beerel/beerserver/templates/index.html 

ENV PRODUCTION=true

EXPOSE 8080/tcp

CMD ["./beerel/beerserver"]
