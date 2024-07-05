# Stage 1

FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . /app

RUN go build -o beerserver main.go

# Stage 2

FROM debian:latest AS runner

WORKDIR /

RUN mkdir -p /beerel 
COPY --from=builder /app/beerserver /beerserver
COPY --from=builder /app/templates/index.html /templates/index.html 
COPY --from=builder /app/dataimport/all.json /dataimport/all.json
COPY --from=builder /app/db/db.schema /db/db.schema

ENV PRODUCTION=true

EXPOSE 8080/tcp

CMD ["./beerserver"]
