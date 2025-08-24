FROM golang:1.23 AS builder

COPY . /src
WORKDIR /src

RUN GOPROXY=https://goproxy.cn make build

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates  \
    netbase \
    && rm -rf /var/lib/apt/lists/ \
    && apt-get autoremove -y && apt-get autoclean -y


# Copy built binaries
COPY --from=builder /src/bin /app
# Copy Swagger UI static files and OpenAPI spec
COPY --from=builder /src/swagger-ui /app/swagger-ui

WORKDIR /app

EXPOSE 8000
EXPOSE 9000
VOLUME /data/conf

CMD ["./api-gateway", "-conf", "/data/conf"]
