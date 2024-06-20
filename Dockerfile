FROM golang:latest as builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .


RUN export CGO_ENABLED=0 && go build -o beerserver server.go

# Stage 2

FROM alpine:latest AS runner

COPY --from=builder /app/beerserver /beerel/beerserver
COPY --from=builder /app/templates/index.html /beerel/beerserver/templates/index.html 

ENV PRODUCTION=true

EXPOSE 8080/tcp

CMD ["./beerel/beerserver"]