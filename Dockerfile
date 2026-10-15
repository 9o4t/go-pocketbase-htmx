# Standalone starter. Build from this directory: docker build .
FROM golang:1.27-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/app .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=build /out/app /usr/local/bin/app
ENV PORT=8090
ENV PB_DATA_DIR=/data
EXPOSE 8090
CMD ["sh", "-c", "exec /usr/local/bin/app serve --http=0.0.0.0:${PORT:-8090}"]
