# Dockerfile.backend
FROM golang:1.21

WORKDIR /app

# Copy binary hasil build
COPY ./app/app /app/app

# Copy env dan logs (jika dibutuhkan)
COPY .env /app/.env
COPY ./logs /app/logs

# Jalankan binary
CMD ["/app/app"]