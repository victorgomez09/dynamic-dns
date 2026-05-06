# --- ETAPA 1: Compilación ---
FROM golang:1.24-alpine AS builder

# Instalamos certificados (necesarios para llamadas HTTPS)
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copiamos archivos de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el resto del código
COPY . .

# Compilamos:
# -w -s: reduce el tamaño del binario eliminando información de depuración
# CGO_ENABLED=0: asegura que el binario sea estático y no dependa de librerías del SO
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o ddns-cloudflare cmd/main.go

# --- ETAPA 2: Imagen Final ---
FROM scratch

# Importamos los certificados desde el builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copiamos el binario compilado
COPY --from=builder /app/ddns-cloudflare /ddns-cloudflare

# Ejecutamos
ENTRYPOINT ["/ddns-cloudflare"]