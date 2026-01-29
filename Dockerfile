# Etapa de construcción
FROM golang:1.25.5-alpine AS builder

# Instalar make y otras dependencias necesarias
RUN apk add --no-cache make gcc musl-dev

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar los archivos necesarios para la compilación
COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY Makefile ./

# Compilar el binario
RUN make build

# Etapa final
FROM alpine:3.18

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar el binario desde la etapa de construcción
COPY --from=builder /app/bin/* ./main

# Dar permisos de ejecución al binario
RUN chmod +x ./main

# Comando por defecto
CMD ["./main"]