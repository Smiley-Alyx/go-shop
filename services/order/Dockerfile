FROM golang:1.22 AS build

WORKDIR /src

COPY go.mod ./

COPY cmd/order ./cmd/order

RUN go build -o /out/order ./cmd/order

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=build /out/order /app/order

EXPOSE 8082

ENV PORT=8082

CMD ["/app/order"]
