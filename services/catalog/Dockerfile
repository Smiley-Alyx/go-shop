FROM golang:1.22 AS build

WORKDIR /src

COPY go.mod ./

COPY cmd/catalog ./cmd/catalog

RUN go build -o /out/catalog ./cmd/catalog

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=build /out/catalog /app/catalog

EXPOSE 8081

ENV PORT=8081

CMD ["/app/catalog"]
