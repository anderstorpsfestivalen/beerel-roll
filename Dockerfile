# Stage 1

FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

ADD . /app

COPY * ./

#RUN apk add --no-cache gcc g++ 

RUN CGO_ENABLED=1 go build -o beerserver main.go

# Stage 2

FROM alpine:latest AS runner

WORKDIR /

RUN mkdir -p /beerel 
COPY --from=builder /app/beerserver /beerel/beerserver
COPY --from=builder /app/templates/index.html /beerel/templates/index.html 
COPY --from=builder /app/dataimport/all.json /beerel/dataimport/all.json
COPY --from=builder /app/db/db.schema /beerel/db/db.schema

ENV PRODUCTION=true

EXPOSE 8080/tcp

CMD ["/beerel/beerserver"]
